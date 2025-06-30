package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var debugMode = os.Getenv("VOICE_NOTIFY_DEBUG") == "true"

// debugLog prints debug messages
func debugLog(format string, args ...interface{}) {
	if !debugMode {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	log.Printf("[DEBUG %s] %s", timestamp, message)
}

// debugMeasureTime measures execution time
func debugMeasureTime(name string) func() {
	if !debugMode {
		return func() {}
	}
	start := time.Now()
	debugLog("Starting %s", name)
	return func() {
		debugLog("Completed %s (took %v)", name, time.Since(start))
	}
}

// debugLogVoiceCommand logs voice command details
func debugLogVoiceCommand(command string, args []string, text string, result error) {
	if !debugMode {
		return
	}

	// Truncate long text for logging
	displayText := text
	if len(displayText) > 50 {
		displayText = displayText[:50] + "..."
	}

	if result != nil {
		debugLog("Voice Command Failed - Command: %s %s, Text: %q, Error: %v",
			command, strings.Join(args, " "), displayText, result)
	} else {
		debugLog("Voice Command Success - Command: %s %s, Text: %q",
			command, strings.Join(args, " "), displayText)
	}
}

// debugLogLanguageDetection logs language detection results
func debugLogLanguageDetection(text string, detectedLang string, source string) {
	if !debugMode {
		return
	}

	// Truncate long text for logging
	displayText := text
	if len(displayText) > 30 {
		displayText = displayText[:30] + "..."
	}

	debugLog("Language Detection - Text: %q, Detected: %s, Source: %s",
		displayText, detectedLang, source)
}

// debugLogEnvironment logs environment variables
func debugLogEnvironment() {
	if !debugMode {
		return
	}

	debugLog("Environment Variables:")
	debugLog("  VOICE_NOTIFY_DEFAULT_VOICE: %s", os.Getenv("VOICE_NOTIFY_DEFAULT_VOICE"))
	debugLog("  VOICE_NOTIFY_DEFAULT_LANGUAGE: %s", os.Getenv("VOICE_NOTIFY_DEFAULT_LANGUAGE"))
	debugLog("  VOICE_NOTIFY_AUTO_DETECT_LANGUAGE: %s", os.Getenv("VOICE_NOTIFY_AUTO_DETECT_LANGUAGE"))
	debugLog("  VOICE_NOTIFY_AUTO_NOTIFY: %s", os.Getenv("VOICE_NOTIFY_AUTO_NOTIFY"))
	debugLog("  VOICE_NOTIFY_MIN_TASK_DURATION: %s", os.Getenv("VOICE_NOTIFY_MIN_TASK_DURATION"))
	debugLog("  VOICE_NOTIFY_QUIET_HOURS: %s", os.Getenv("VOICE_NOTIFY_QUIET_HOURS"))
	debugLog("  VOICE_NOTIFY_DEBUG: %s", os.Getenv("VOICE_NOTIFY_DEBUG"))
}

// debugLogRequest logs MCP request details
func debugLogRequest(method string, params interface{}) {
	if !debugMode {
		return
	}

	debugLog("MCP Request - Method: %s, Params: %+v", method, params)
}

// debugLogRateLimit logs rate limit information
func debugLogRateLimit(allowed bool, reason string) {
	if !debugMode {
		return
	}

	if allowed {
		debugLog("Rate Limit - Allowed: %s", reason)
	} else {
		debugLog("Rate Limit - Denied: %s", reason)
	}
}

// debugLogVoiceSelection logs voice selection process
func debugLogVoiceSelection(stage string, voice string, reason string) {
	if !debugMode {
		return
	}

	debugLog("Voice Selection - Stage: %s, Voice: %s, Reason: %s", stage, voice, reason)
}
