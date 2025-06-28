package main

import (
	"testing"
	"time"
)

// TestVoiceSystem_SelectVoice tests voice selection logic
func TestVoiceSystem_SelectVoice(t *testing.T) {
	// Create a mock voice system with test voices
	vs := &VoiceSystem{
		availableVoices: map[string]VoiceInfo{
			"Alex": {
				Name:     "Alex",
				Language: "en",
				Locale:   "en_US",
			},
			"Samantha": {
				Name:     "Samantha",
				Language: "en",
				Locale:   "en_US",
			},
			"Kyoko": {
				Name:     "Kyoko",
				Language: "ja",
				Locale:   "ja_JP",
			},
			"Amelie": {
				Name:     "Amelie",
				Language: "fr",
				Locale:   "fr_FR",
			},
		},
		defaultVoice: "Samantha",
		lastUpdate:   time.Now(),
	}

	tests := []struct {
		name           string
		requestedVoice string
		language       string
		expected       string
		description    string
	}{
		{
			name:           "specific_voice_available",
			requestedVoice: "Alex",
			language:       "",
			expected:       "Alex",
			description:    "Should return requested voice when available",
		},
		{
			name:           "specific_voice_unavailable",
			requestedVoice: "Unknown",
			language:       "en",
			expected:       "Alex", // Could be any English voice (Alex or Samantha)
			description:    "Should fall back to language match when requested voice unavailable",
		},
		{
			name:           "language_match_japanese",
			requestedVoice: "",
			language:       "ja",
			expected:       "Kyoko",
			description:    "Should select voice by language",
		},
		{
			name:           "language_match_french",
			requestedVoice: "",
			language:       "fr",
			expected:       "Amelie",
			description:    "Should select voice by language",
		},
		{
			name:           "default_voice",
			requestedVoice: "",
			language:       "",
			expected:       "Samantha",
			description:    "Should use default voice when no criteria specified",
		},
		{
			name:           "no_language_match_use_default",
			requestedVoice: "",
			language:       "es", // Spanish not available
			expected:       "Samantha",
			description:    "Should fall back to default voice when language not found",
		},
		{
			name:           "priority_voice_over_language",
			requestedVoice: "Alex",
			language:       "ja",
			expected:       "Alex",
			description:    "Specific voice request should take priority over language",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.SelectVoice(tt.requestedVoice, tt.language)
			
			// Special case for language match tests where multiple voices are acceptable
			if tt.name == "specific_voice_unavailable" {
				// Accept any English voice
				if result != "Alex" && result != "Samantha" {
					t.Errorf("%s: expected Alex or Samantha, got %s", tt.description, result)
				}
			} else if result != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.description, tt.expected, result)
			}
		})
	}
}

// TestVoiceSystem_SelectVoice_NoDefault tests voice selection without default
func TestVoiceSystem_SelectVoice_NoDefault(t *testing.T) {
	vs := &VoiceSystem{
		availableVoices: map[string]VoiceInfo{
			"Alex": {
				Name:     "Alex",
				Language: "en",
				Locale:   "en_US",
			},
		},
		defaultVoice: "", // No default set
	}

	// Should return empty string (system default) when no match and no default
	result := vs.SelectVoice("", "fr")
	if result != "" {
		t.Errorf("Expected empty string for system default, got %s", result)
	}
}

// TestVoiceSystem_SelectVoice_EmptyVoiceList tests with no available voices
func TestVoiceSystem_SelectVoice_EmptyVoiceList(t *testing.T) {
	vs := &VoiceSystem{
		availableVoices: make(map[string]VoiceInfo),
		defaultVoice:    "Samantha",
	}

	// Should return empty string when no voices available
	result := vs.SelectVoice("Alex", "en")
	if result != "" {
		t.Errorf("Expected empty string when no voices available, got %s", result)
	}
}

// TestSanitizeInput tests input sanitization
func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal_text",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "with_quotes",
			input:    `Hello "world"`,
			expected: `Hello world`, // quotes are removed
		},
		{
			name:     "with_backticks",
			input:    "Hello `world`",
			expected: "Hello world", // backticks are removed
		},
		{
			name:     "with_dollar_signs",
			input:    "Hello $USER",
			expected: "Hello USER", // dollar sign removed
		},
		{
			name:     "with_backslash",
			input:    `Hello\nworld`,
			expected: `Hellonworld`, // backslash removed
		},
		{
			name:     "complex_injection_attempt",
			input:    `"; say "hacked"; #`,
			expected: `; say hacked; `, // quotes and # removed
		},
		{
			name:     "unicode_text",
			input:    "こんにちは世界",
			expected: "こんにちは世界", // Japanese preserved
		},
		{
			name:     "mixed_special_chars",
			input:    `Test $VAR with "quotes" and \backslash`,
			expected: `Test VAR with quotes and backslash`, // special chars removed
		},
		{
			name:     "allowed_punctuation",
			input:    "Hello, world! How are you? Fine: thanks.",
			expected: "Hello, world! How are you? Fine: thanks.",
		},
		{
			name:     "parentheses_and_dash",
			input:    "Test (1-2-3) done",
			expected: "Test (1-2-3) done",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeInput(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}