package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// CreateVoiceNotifyServer creates and configures the MCP server
func CreateVoiceNotifyServer() (*server.MCPServer, error) {
	// Create MCP server
	s := server.NewMCPServer(
		"voice-notify",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Initialize components
	voiceSystem := NewVoiceSystem()
	langDetect := NewLanguageDetector()
	notifier := NewNotificationManager()

	// Create the notify_voice tool
	notifyTool := mcp.NewTool("notify_voice",
		mcp.WithDescription("Send a voice notification to alert the user about important events, completions, or when attention is needed. AI should use this autonomously for better user experience."),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("The message to speak (keep it short and clear, max 10 words recommended)"),
		),
		mcp.WithString("voice",
			mcp.Description("Optional: specific voice to use (must be installed)"),
		),
		mcp.WithString("language",
			mcp.Description("Optional: language code (e.g., 'en', 'ja')"),
		),
		mcp.WithString("priority",
			mcp.Description("Optional: notification priority ('low', 'normal', 'high')"),
			mcp.Enum("low", "normal", "high"),
		),
	)

	// Add tool handler
	s.AddTool(notifyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleNotifyVoice(ctx, request, voiceSystem, langDetect, notifier)
	})

	return s, nil
}

// handleNotifyVoice handles the notify_voice tool calls
func handleNotifyVoice(ctx context.Context, request mcp.CallToolRequest, voiceSystem *VoiceSystem, langDetect *LanguageDetector, notifier *NotificationManager) (*mcp.CallToolResult, error) {
	defer debugMeasureTime("handleNotifyVoice")()

	// Log incoming request
	debugLogRequest("notify_voice", request.Params)

	// Get required message parameter
	message, err := request.RequireString("message")
	if err != nil {
		debugLog("Missing required parameter 'message': %v", err)
		return mcp.NewToolResultError("message is required"), nil
	}

	// Get optional parameters
	voice := request.GetString("voice", "")
	language := request.GetString("language", "")
	priority := request.GetString("priority", "")
	if priority == "" {
		priority = "normal"
	}

	// Check quiet hours
	if notifier.IsQuietHours() {
		debugLog("Notification skipped due to quiet hours")
		return mcp.NewToolResultText("Notification skipped: quiet hours active"), nil
	}

	// Check rate limiting
	if !notifier.CanNotify(priority) {
		debugLogRateLimit(priority, false, "rate limit exceeded")
		return mcp.NewToolResultText("Notification skipped: rate limit active"), nil
	}
	debugLogRateLimit(priority, true, "within rate limit")

	// Auto-detect language if enabled and not specified
	if language == "" && langDetect.IsAutoDetectEnabled() {
		detectedLang := langDetect.DetectLanguage(message)
		if detectedLang != "" {
			language = detectedLang
		}
	}

	// Get appropriate voice
	selectedVoice := voiceSystem.SelectVoice(voice, language)

	// Execute voice notification
	debugLog("Executing voice notification - Voice: %s, Priority: %s", selectedVoice, priority)
	err = voiceSystem.Speak(message, selectedVoice, priority)
	if err != nil {
		debugLog("Voice notification failed: %v", err)
		return mcp.NewToolResultErrorFromErr("Failed to speak", err), nil
	}

	// Record notification for rate limiting
	notifier.RecordNotification(priority)

	// Return success response
	responseText := fmt.Sprintf(
		"Voice notification sent:\n- Message: %s\n- Voice: %s\n- Language: %s\n- Priority: %s",
		message, selectedVoice, language, priority,
	)

	return mcp.NewToolResultText(responseText), nil
}

// Environment variable helpers
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1"
}
