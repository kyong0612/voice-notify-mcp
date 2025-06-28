package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// VoiceSystem manages voice synthesis using macOS 'say' command
type VoiceSystem struct {
	availableVoices map[string]VoiceInfo
	defaultVoice    string
	mu              sync.RWMutex
	lastUpdate      time.Time
}

// VoiceInfo contains information about a voice
type VoiceInfo struct {
	Name     string
	Language string
	Locale   string
}

// NewVoiceSystem creates a new voice system instance
func NewVoiceSystem() *VoiceSystem {
	vs := &VoiceSystem{
		availableVoices: make(map[string]VoiceInfo),
		defaultVoice:    getEnv("VOICE_NOTIFY_DEFAULT_VOICE", ""),
	}
	
	// Load available voices
	vs.refreshVoices()
	
	return vs
}

// refreshVoices updates the list of available voices
func (vs *VoiceSystem) refreshVoices() error {
	defer debugMeasureTime("refreshVoices")()
	
	vs.mu.Lock()
	defer vs.mu.Unlock()

	// Run 'say -v ?' command to get available voices
	cmd := exec.Command("say", "-v", "?")
	debugLogVoiceCommand("say", []string{"-v", "?"}, "", nil)
	output, err := cmd.Output()
	if err != nil {
		debugLog("Failed to get voice list: %v", err)
		return fmt.Errorf("failed to get voices: %w", err)
	}

	// Parse output
	// Format: "Name             Language         Sample"
	vs.availableVoices = make(map[string]VoiceInfo)
	debugLog("Parsing voice list output (%d bytes)", len(output))
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}

		// Parse voice info (name and language are separated by spaces)
		// Example: "Alex                en_US    # Most people recognize me by my voice."
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			locale := parts[1]
			
			// Extract language code from locale
			lang := strings.Split(locale, "_")[0]
			
			vs.availableVoices[name] = VoiceInfo{
				Name:     name,
				Language: lang,
				Locale:   locale,
			}
		}
	}

	vs.lastUpdate = time.Now()
	debugLog("Loaded %d voices", len(vs.availableVoices))
	return nil
}

// SelectVoice selects the appropriate voice based on preferences
func (vs *VoiceSystem) SelectVoice(requestedVoice, language string) string {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	// 1. If specific voice is requested and available, use it
	if requestedVoice != "" {
		if _, exists := vs.availableVoices[requestedVoice]; exists {
			debugLogVoiceSelection(requestedVoice, language, requestedVoice, false)
			return requestedVoice
		}
		debugLog("Requested voice '%s' not available", requestedVoice)
	}

	// 2. If language is specified, find a voice for that language
	if language != "" {
		for name, info := range vs.availableVoices {
			if info.Language == language {
				debugLogVoiceSelection(requestedVoice, language, name, true)
				return name
			}
		}
		debugLog("No voice found for language '%s'", language)
	}

	// 3. Use default voice if set and available
	if vs.defaultVoice != "" {
		if _, exists := vs.availableVoices[vs.defaultVoice]; exists {
			debugLogVoiceSelection(requestedVoice, language, vs.defaultVoice, true)
			return vs.defaultVoice
		}
		debugLog("Default voice '%s' not available", vs.defaultVoice)
	}

	// 4. Use system default (empty string means use system default)
	debugLogVoiceSelection(requestedVoice, language, "", true)
	return ""
}

// Speak executes the say command with the given message and voice
func (vs *VoiceSystem) Speak(message, voice, priority string) error {
	// Sanitize input to prevent command injection
	message = sanitizeInput(message)
	
	// Build command arguments
	args := []string{}
	
	// Add voice if specified
	if voice != "" {
		args = append(args, "-v", voice)
	}
	
	// Adjust rate based on priority
	switch priority {
	case "high":
		args = append(args, "-r", "200") // Faster speech
	case "low":
		args = append(args, "-r", "150") // Slower speech
	default:
		// Normal rate (default)
	}
	
	// Add the message
	args = append(args, message)
	
	// Execute command
	cmd := exec.Command("say", args...)
	
	// Capture both stdout and stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	debugLogVoiceCommand("say", args, "", nil)
	err := cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("say command failed: %v, stderr: %s", err, stderr.String())
		debugLogVoiceCommand("say", args, "", fmt.Errorf("%s", errMsg))
		return fmt.Errorf("say command failed: %w, stderr: %s", err, stderr.String())
	}
	
	return nil
}

// GetAvailableVoices returns a list of available voices
func (vs *VoiceSystem) GetAvailableVoices() []VoiceInfo {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	// Refresh if data is older than 5 minutes
	if time.Since(vs.lastUpdate) > 5*time.Minute {
		go vs.refreshVoices()
	}

	voices := make([]VoiceInfo, 0, len(vs.availableVoices))
	for _, voice := range vs.availableVoices {
		voices = append(voices, voice)
	}
	
	return voices
}

// sanitizeInput sanitizes the input to prevent command injection
func sanitizeInput(input string) string {
	// Remove any characters that could be used for command injection
	// Allow only safe characters
	var sanitized strings.Builder
	
	for _, r := range input {
		// Allow alphanumeric, common punctuation, and spaces
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == ' ' || r == '.' ||
			r == ',' || r == '!' || r == '?' || r == '-' ||
			r == ':' || r == ';' || r == '(' || r == ')' ||
			(r >= 0x3040 && r <= 0x309F) || // Hiragana
			(r >= 0x30A0 && r <= 0x30FF) || // Katakana
			(r >= 0x4E00 && r <= 0x9FAF) || // CJK Unified Ideographs
			(r >= 0xAC00 && r <= 0xD7AF) {  // Hangul
			sanitized.WriteRune(r)
		}
	}
	
	return sanitized.String()
}