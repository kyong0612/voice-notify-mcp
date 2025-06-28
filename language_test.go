package main

import (
	"testing"
)

// TestLanguageDetector_DetectLanguage tests language detection
func TestLanguageDetector_DetectLanguage(t *testing.T) {
	ld := &LanguageDetector{
		autoDetect:      true,
		defaultLanguage: "en",
	}

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		// English tests
		{
			name:     "english_simple",
			text:     "Hello world",
			expected: "en",
		},
		{
			name:     "english_sentence",
			text:     "The quick brown fox jumps over the lazy dog",
			expected: "en",
		},
		
		// Japanese tests
		{
			name:     "japanese_hiragana",
			text:     "こんにちは",
			expected: "ja",
		},
		{
			name:     "japanese_katakana",
			text:     "コンピューター",
			expected: "ja",
		},
		{
			name:     "japanese_mixed",
			text:     "今日は良い天気ですね",
			expected: "ja",
		},
		{
			name:     "japanese_with_english",
			text:     "Hello こんにちは",
			expected: "ja", // Japanese takes priority
		},

		// Chinese tests
		{
			name:     "chinese_simple",
			text:     "你好世界",
			expected: "zh",
		},
		{
			name:     "chinese_traditional",
			text:     "歡迎使用語音通知",
			expected: "zh",
		},

		// Korean tests
		{
			name:     "korean_hangul",
			text:     "안녕하세요",
			expected: "ko",
		},
		{
			name:     "korean_sentence",
			text:     "오늘 날씨가 좋네요",
			expected: "ko",
		},

		// Russian tests
		{
			name:     "russian_cyrillic",
			text:     "Привет мир",
			expected: "ru",
		},
		{
			name:     "russian_sentence",
			text:     "Добро пожаловать в систему голосовых уведомлений",
			expected: "ru",
		},

		// French tests
		{
			name:     "french_articles",
			text:     "Le chat est sur la table",
			expected: "fr",
		},
		{
			name:     "french_accents",
			text:     "Félicitations! Vous avez réussi",
			expected: "fr",
		},

		// Spanish tests
		{
			name:     "spanish_articles",
			text:     "El gato está en la mesa",
			expected: "es",
		},
		{
			name:     "spanish_special",
			text:     "¿Cómo estás? ¡Muy bien!",
			expected: "es",
		},

		// German tests
		{
			name:     "german_articles",
			text:     "Der Hund ist in dem Garten",
			expected: "de",
		},
		{
			name:     "german_umlauts",
			text:     "Schöne Grüße aus München",
			expected: "de",
		},

		// Italian tests
		{
			name:     "italian_articles",
			text:     "Il cane è nel giardino",
			expected: "it",
		},

		// Portuguese tests
		{
			name:     "portuguese_articles",
			text:     "O gato está na mesa",
			expected: "pt",
		},
		{
			name:     "portuguese_special",
			text:     "Atenção! A operação foi concluída",
			expected: "pt",
		},

		// Edge cases
		{
			name:     "empty_string",
			text:     "",
			expected: "en", // default
		},
		{
			name:     "numbers_only",
			text:     "123456789",
			expected: "en", // default
		},
		{
			name:     "mixed_scripts_english_dominant",
			text:     "Hello world 123 test",
			expected: "en",
		},
		{
			name:     "unknown_script",
			text:     "☺️🎉🎊✨",
			expected: "en", // default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ld.DetectLanguage(tt.text)
			if result != tt.expected {
				t.Errorf("DetectLanguage(%q) = %q, want %q", tt.text, result, tt.expected)
			}
		})
	}
}

// TestLanguageDetector_AutoDetectDisabled tests when auto-detect is disabled
func TestLanguageDetector_AutoDetectDisabled(t *testing.T) {
	ld := &LanguageDetector{
		autoDetect:      false,
		defaultLanguage: "fr",
	}

	// Should always return default when auto-detect is disabled
	tests := []struct {
		text string
	}{
		{"Hello world"},
		{"こんにちは"},
		{"Bonjour"},
	}

	for _, tt := range tests {
		result := ld.DetectLanguage(tt.text)
		if result != "fr" {
			t.Errorf("With auto-detect disabled, expected 'fr', got %q for text %q", result, tt.text)
		}
	}
}

// TestContainsAny tests the helper function
func TestContainsAny(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		patterns  []string
		expected  bool
	}{
		{
			name:      "contains_pattern",
			text:      "hello world",
			patterns:  []string{"world", "test"},
			expected:  true,
		},
		{
			name:      "contains_multiple",
			text:      "the quick brown fox",
			patterns:  []string{"quick", "brown"},
			expected:  true,
		},
		{
			name:      "contains_none",
			text:      "hello world",
			patterns:  []string{"foo", "bar"},
			expected:  false,
		},
		{
			name:      "empty_patterns",
			text:      "hello world",
			patterns:  []string{},
			expected:  false,
		},
		{
			name:      "empty_text",
			text:      "",
			patterns:  []string{"test"},
			expected:  false,
		},
		{
			name:      "case_sensitive",
			text:      "Hello World",
			patterns:  []string{"hello"}, // lowercase
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsAny(tt.text, tt.patterns)
			if result != tt.expected {
				t.Errorf("containsAny(%q, %v) = %v, want %v", tt.text, tt.patterns, result, tt.expected)
			}
		})
	}
}

// TestCharacterTypeDetection tests individual character detection functions
func TestCharacterTypeDetection(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		isLatin  bool
		isJapanese bool
		isChinese bool
		isKorean bool
		isCyrillic bool
		isArabic bool
		isHebrew bool
	}{
		{
			name:     "latin_lowercase",
			char:     'a',
			isLatin:  true,
		},
		{
			name:     "latin_uppercase",
			char:     'Z',
			isLatin:  true,
		},
		{
			name:       "japanese_hiragana",
			char:       'あ',
			isJapanese: true,
		},
		{
			name:       "japanese_katakana",
			char:       'ア',
			isJapanese: true,
		},
		{
			name:      "chinese_hanzi",
			char:      '中',
			isChinese: true,
		},
		{
			name:     "korean_hangul",
			char:     '한',
			isKorean: true,
		},
		{
			name:       "cyrillic",
			char:       'Д',
			isCyrillic: true,
		},
		{
			name:     "arabic",
			char:     'م',
			isArabic: true,
		},
		{
			name:     "hebrew",
			char:     'ש',
			isHebrew: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLatin(tt.char); got != tt.isLatin {
				t.Errorf("isLatin(%c) = %v, want %v", tt.char, got, tt.isLatin)
			}
			if got := isJapanese(tt.char); got != tt.isJapanese {
				t.Errorf("isJapanese(%c) = %v, want %v", tt.char, got, tt.isJapanese)
			}
			if got := isChinese(tt.char); got != tt.isChinese {
				t.Errorf("isChinese(%c) = %v, want %v", tt.char, got, tt.isChinese)
			}
			if got := isKorean(tt.char); got != tt.isKorean {
				t.Errorf("isKorean(%c) = %v, want %v", tt.char, got, tt.isKorean)
			}
			if got := isCyrillic(tt.char); got != tt.isCyrillic {
				t.Errorf("isCyrillic(%c) = %v, want %v", tt.char, got, tt.isCyrillic)
			}
			if got := isArabic(tt.char); got != tt.isArabic {
				t.Errorf("isArabic(%c) = %v, want %v", tt.char, got, tt.isArabic)
			}
			if got := isHebrew(tt.char); got != tt.isHebrew {
				t.Errorf("isHebrew(%c) = %v, want %v", tt.char, got, tt.isHebrew)
			}
		})
	}
}