package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// NotificationManager manages notification logic and rate limiting
type NotificationManager struct {
	autoNotify      bool
	minTaskDuration time.Duration
	quietHours      *QuietHours
	lastNotif       map[string]time.Time
	mu              sync.RWMutex
}

// QuietHours represents a time range for quiet hours
type QuietHours struct {
	Start time.Time
	End   time.Time
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager() *NotificationManager {
	nm := &NotificationManager{
		autoNotify:      getEnvBool("VOICE_NOTIFY_AUTO_NOTIFY", true),
		minTaskDuration: parseMinTaskDuration(),
		lastNotif:       make(map[string]time.Time),
	}

	// Parse quiet hours
	if quietHoursStr := getEnv("VOICE_NOTIFY_QUIET_HOURS", ""); quietHoursStr != "" {
		nm.quietHours = parseQuietHours(quietHoursStr)
		if nm.quietHours != nil {
			debugLog("Quiet hours configured: %s to %s",
				nm.quietHours.Start.Format("15:04"),
				nm.quietHours.End.Format("15:04"))
		}
	}

	debugLog("NotificationManager initialized - AutoNotify: %v, MinTaskDuration: %v",
		nm.autoNotify, nm.minTaskDuration)

	return nm
}

// IsAutoNotifyEnabled returns whether auto-notification is enabled
func (nm *NotificationManager) IsAutoNotifyEnabled() bool {
	return nm.autoNotify
}

// ShouldNotify determines if a notification should be sent based on task duration
func (nm *NotificationManager) ShouldNotify(taskDuration time.Duration) bool {
	if !nm.autoNotify {
		debugLog("Auto-notify disabled, skipping notification")
		return false
	}
	return taskDuration >= nm.minTaskDuration
}

// IsQuietHours checks if current time is within quiet hours
func (nm *NotificationManager) IsQuietHours() bool {
	if nm.quietHours == nil {
		return false
	}

	now := time.Now()
	currentTime := time.Date(0, 1, 1, now.Hour(), now.Minute(), 0, 0, time.Local)

	// Handle quiet hours that span midnight
	if nm.quietHours.End.Before(nm.quietHours.Start) {
		// Quiet hours span midnight (e.g., 22:00 - 07:00)
		isQuiet := currentTime.After(nm.quietHours.Start) || currentTime.Before(nm.quietHours.End)
		debugLog("Quiet hours check (spans midnight): Current=%s, Start=%s, End=%s, IsQuiet=%v",
			currentTime.Format("15:04"), nm.quietHours.Start.Format("15:04"),
			nm.quietHours.End.Format("15:04"), isQuiet)
		return isQuiet
	}

	// Normal quiet hours (e.g., 23:00 - 06:00)
	isQuiet := currentTime.After(nm.quietHours.Start) && currentTime.Before(nm.quietHours.End)
	debugLog("Quiet hours check: Current=%s, Start=%s, End=%s, IsQuiet=%v",
		currentTime.Format("15:04"), nm.quietHours.Start.Format("15:04"),
		nm.quietHours.End.Format("15:04"), isQuiet)
	return isQuiet
}

// RecordNotification records a notification for rate limiting
func (nm *NotificationManager) RecordNotification(priority string) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	nm.lastNotif[priority] = time.Now()
}

// CanNotify checks if enough time has passed since the last notification of this type
func (nm *NotificationManager) CanNotify(priority string) bool {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	lastTime, exists := nm.lastNotif[priority]
	if !exists {
		debugLog("No previous notification for priority '%s', allowing notification", priority)
		return true
	}

	// Rate limits based on priority
	var minInterval time.Duration
	switch priority {
	case "high":
		minInterval = 10 * time.Second
	case "normal":
		minInterval = 30 * time.Second
	case "low":
		minInterval = 60 * time.Second
	default:
		minInterval = 30 * time.Second
	}

	elapsed := time.Since(lastTime)
	canNotify := elapsed >= minInterval
	debugLog("Rate limit check - Priority: %s, MinInterval: %v, Elapsed: %v, CanNotify: %v",
		priority, minInterval, elapsed, canNotify)

	return canNotify
}

// parseMinTaskDuration parses the minimum task duration from environment
func parseMinTaskDuration() time.Duration {
	durationStr := getEnv("VOICE_NOTIFY_MIN_TASK_DURATION", "3")
	seconds, err := strconv.Atoi(durationStr)
	if err != nil || seconds < 0 {
		return 3 * time.Second
	}
	return time.Duration(seconds) * time.Second
}

// parseQuietHours parses quiet hours from a string format like "22:00-07:00"
func parseQuietHours(quietHoursStr string) *QuietHours {
	parts := strings.Split(quietHoursStr, "-")
	if len(parts) != 2 {
		return nil
	}

	start, err := parseTime(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil
	}

	end, err := parseTime(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil
	}

	return &QuietHours{
		Start: start,
		End:   end,
	}
}

// parseTime parses a time string in HH:MM format
func parseTime(timeStr string) (time.Time, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, err
	}
	if hour < 0 || hour > 23 {
		return time.Time{}, fmt.Errorf("invalid hour: %d", hour)
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, err
	}
	if minute < 0 || minute > 59 {
		return time.Time{}, fmt.Errorf("invalid minute: %d", minute)
	}

	return time.Date(0, 1, 1, hour, minute, 0, 0, time.Local), nil
}
