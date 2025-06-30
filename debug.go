package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// Debug configuration
var (
	debugMode   bool
	debugLogger *log.Logger
)

func init() {
	// Check if debug mode is enabled
	debugMode = getEnvBool("VOICE_NOTIFY_DEBUG", false)

	if debugMode {
		// Set up debug logger with more detailed format
		debugLogger = log.New(os.Stderr, "[DEBUG] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	}
}

// debugLog logs a message if debug mode is enabled
func debugLog(format string, args ...interface{}) {
	if !debugMode || debugLogger == nil {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(1)
	if ok {
		// Extract just the filename
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]

		// Prepend file:line to the message
		format = fmt.Sprintf("[%s:%d] %s", file, line, format)
	}

	debugLogger.Printf(format, args...)
}

// debugLogRequest logs MCP request details
func debugLogRequest(requestType string, data interface{}) {
	if !debugMode {
		return
	}

	debugLog("MCP Request - Type: %s, Data: %+v", requestType, data)
}

// debugLogResponse logs MCP response details
func debugLogResponse(responseType string, data interface{}, err error) {
	if !debugMode {
		return
	}

	if err != nil {
		debugLog("MCP Response - Type: %s, Error: %v", responseType, err)
	} else {
		debugLog("MCP Response - Type: %s, Data: %+v", responseType, data)
	}
}

// debugLogVoiceCommand logs voice command details
func debugLogVoiceCommand(command string, args []string, output string, err error) {
	if !debugMode {
		return
	}

	debugLog("Voice Command - Cmd: %s, Args: %v", command, args)
	if output != "" {
		debugLog("Voice Output: %s", output)
	}
	if err != nil {
		debugLog("Voice Error: %v", err)
	}
}

// debugLogEnvironment logs environment configuration at startup
func debugLogEnvironment() {
	if !debugMode {
		return
	}

	debugLog("=== Environment Configuration ===")
	debugLog("VOICE_NOTIFY_DEBUG: %v", debugMode)
	debugLog("VOICE_NOTIFY_DEFAULT_VOICE: %s", getEnv("VOICE_NOTIFY_DEFAULT_VOICE", ""))
	debugLog("VOICE_NOTIFY_DEFAULT_LANGUAGE: %s", getEnv("VOICE_NOTIFY_DEFAULT_LANGUAGE", ""))
	debugLog("VOICE_NOTIFY_AUTO_DETECT_LANGUAGE: %v", getEnvBool("VOICE_NOTIFY_AUTO_DETECT_LANGUAGE", true))
	debugLog("VOICE_NOTIFY_AUTO_NOTIFY: %v", getEnvBool("VOICE_NOTIFY_AUTO_NOTIFY", true))
	debugLog("VOICE_NOTIFY_MIN_TASK_DURATION: %s", getEnv("VOICE_NOTIFY_MIN_TASK_DURATION", "3"))
	debugLog("VOICE_NOTIFY_QUIET_HOURS: %s", getEnv("VOICE_NOTIFY_QUIET_HOURS", ""))
	debugLog("================================")
}

// debugMeasureTime measures and logs execution time
func debugMeasureTime(operation string) func() {
	if !debugMode {
		return func() {}
	}

	start := time.Now()
	debugLog("Starting operation: %s", operation)

	return func() {
		duration := time.Since(start)
		debugLog("Completed operation: %s (took %v)", operation, duration)
	}
}

// debugLogLanguageDetection logs language detection details
func debugLogLanguageDetection(text string, detectedLang string, confidence float64) {
	if !debugMode {
		return
	}

	debugLog("Language Detection - Text: '%s', Detected: %s, Confidence: %.2f",
		truncateString(text, 50), detectedLang, confidence)
}

// debugLogVoiceSelection logs voice selection process
func debugLogVoiceSelection(requestedVoice, language, selectedVoice string, fallback bool) {
	if !debugMode {
		return
	}

	debugLog("Voice Selection - Requested: %s, Language: %s, Selected: %s, Fallback: %v",
		requestedVoice, language, selectedVoice, fallback)
}

// debugLogRateLimit logs rate limiting decisions
func debugLogRateLimit(priority string, allowed bool, reason string) {
	if !debugMode {
		return
	}

	debugLog("Rate Limit Check - Priority: %s, Allowed: %v, Reason: %s",
		priority, allowed, reason)
}

// truncateString truncates a string for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
