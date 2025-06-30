//go:build integration
// +build integration

package main

import (
	"testing"
	"time"
)

// TestIntegration_ComponentCreation tests component creation
func TestIntegration_ComponentCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test VoiceSystem creation
	voiceSystem := NewVoiceSystem()
	if voiceSystem == nil {
		t.Fatal("VoiceSystem is nil")
	}

	// Test LanguageDetector creation
	langDetect := NewLanguageDetector()
	if langDetect == nil {
		t.Fatal("LanguageDetector is nil")
	}

	// Test NotificationManager creation
	notifier := NewNotificationManager()
	if notifier == nil {
		t.Fatal("NotificationManager is nil")
	}

	// Test server creation
	server, err := CreateVoiceNotifyServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	if server == nil {
		t.Fatal("Server is nil")
	}
}

// TestIntegration_RateLimiting tests rate limiting behavior
func TestIntegration_RateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	notifier := NewNotificationManager()

	// Test rapid notifications
	priorities := []string{"high", "normal", "low"}

	for _, priority := range priorities {
		t.Run("rate_limit_"+priority, func(t *testing.T) {
			// First notification should succeed
			if !notifier.CanNotify(priority) {
				t.Errorf("First %s notification should be allowed", priority)
			}
			notifier.RecordNotification(priority)

			// Immediate second should fail
			if notifier.CanNotify(priority) {
				t.Errorf("Immediate second %s notification should be blocked", priority)
			}

			// Wait a bit (not full duration)
			time.Sleep(2 * time.Second)

			// Should still be blocked
			if notifier.CanNotify(priority) {
				t.Errorf("%s notification should still be blocked", priority)
			}
		})
	}
}

// TestIntegration_LanguageDetection tests language detection with real text
func TestIntegration_LanguageDetection(t *testing.T) {
	detector := NewLanguageDetector()

	realTexts := []struct {
		text     string
		expected string
	}{
		{
			text:     "Hello! The build is complete. All tests passed successfully.",
			expected: "en",
		},
		{
			text:     "ビルドが完了しました。すべてのテストが正常に完了しました。",
			expected: "ja",
		},
		{
			text:     "La construction est terminée. Tous les tests ont réussi.",
			expected: "fr",
		},
		{
			text:     "Der Build ist abgeschlossen. Alle Tests waren erfolgreich.",
			expected: "de",
		},
		{
			text:     "La compilación está completa. Todas las pruebas pasaron.",
			expected: "es",
		},
		{
			text:     "你好！构建完成了。所有测试都通过了。",
			expected: "zh",
		},
		{
			text:     "안녕하세요! 빌드가 완료되었습니다. 모든 테스트가 통과했습니다.",
			expected: "ko",
		},
	}

	for _, tt := range realTexts {
		t.Run(tt.expected+"_text", func(t *testing.T) {
			result := detector.DetectLanguage(tt.text)
			if result != tt.expected {
				t.Errorf("DetectLanguage(%q) = %q, want %q", tt.text, result, tt.expected)
			}
		})
	}
}

// MockVoiceSystem is a mock implementation for testing
type MockVoiceSystem struct {
	SpeakCalls []SpeakCall
	Voices     map[string]VoiceInfo
}

type SpeakCall struct {
	Message  string
	Voice    string
	Priority string
}

func (m *MockVoiceSystem) Speak(message, voice, priority string) error {
	m.SpeakCalls = append(m.SpeakCalls, SpeakCall{
		Message:  message,
		Voice:    voice,
		Priority: priority,
	})
	return nil
}

func (m *MockVoiceSystem) SelectVoice(requestedVoice, language string) string {
	if requestedVoice != "" && m.Voices != nil {
		if _, exists := m.Voices[requestedVoice]; exists {
			return requestedVoice
		}
	}

	if language != "" && m.Voices != nil {
		for name, info := range m.Voices {
			if info.Language == language {
				return name
			}
		}
	}

	return "default"
}

// TestWithMockVoice tests voice notification with a mock
func TestWithMockVoice(t *testing.T) {
	mock := &MockVoiceSystem{
		SpeakCalls: make([]SpeakCall, 0),
		Voices: map[string]VoiceInfo{
			"Alex":   {Name: "Alex", Language: "en"},
			"Kyoko":  {Name: "Kyoko", Language: "ja"},
			"Amelie": {Name: "Amelie", Language: "fr"},
		},
	}

	// Test speak calls
	testCases := []struct {
		message  string
		voice    string
		priority string
	}{
		{"Hello world", "Alex", "normal"},
		{"こんにちは", "Kyoko", "high"},
		{"Bonjour", "Amelie", "low"},
	}

	for _, tc := range testCases {
		err := mock.Speak(tc.message, tc.voice, tc.priority)
		if err != nil {
			t.Errorf("Mock speak failed: %v", err)
		}
	}

	// Verify calls
	if len(mock.SpeakCalls) != len(testCases) {
		t.Errorf("Expected %d speak calls, got %d", len(testCases), len(mock.SpeakCalls))
	}

	for i, call := range mock.SpeakCalls {
		if call.Message != testCases[i].message {
			t.Errorf("Call %d: expected message %q, got %q", i, testCases[i].message, call.Message)
		}
		if call.Voice != testCases[i].voice {
			t.Errorf("Call %d: expected voice %q, got %q", i, testCases[i].voice, call.Voice)
		}
		if call.Priority != testCases[i].priority {
			t.Errorf("Call %d: expected priority %q, got %q", i, testCases[i].priority, call.Priority)
		}
	}
}
