package main

import (
	"strings"
	"unicode"
)

// LanguageDetector detects language from text
type LanguageDetector struct {
	autoDetect      bool
	defaultLanguage string
}

// NewLanguageDetector creates a new language detector
func NewLanguageDetector() *LanguageDetector {
	return &LanguageDetector{
		autoDetect:      getEnvBool("VOICE_NOTIFY_AUTO_DETECT_LANGUAGE", true),
		defaultLanguage: getEnv("VOICE_NOTIFY_DEFAULT_LANGUAGE", "en"),
	}
}

// IsAutoDetectEnabled returns whether auto-detection is enabled
func (ld *LanguageDetector) IsAutoDetectEnabled() bool {
	return ld.autoDetect
}

// DetectLanguage detects the language of the given text
func (ld *LanguageDetector) DetectLanguage(text string) string {
	if !ld.autoDetect {
		return ld.defaultLanguage
	}

	// Count character types
	var (
		latin    int
		japanese int
		chinese  int
		korean   int
		cyrillic int
		arabic   int
		hebrew   int
	)

	for _, r := range text {
		switch {
		case isLatin(r):
			latin++
		case isJapanese(r):
			japanese++
		case isChinese(r):
			chinese++
		case isKorean(r):
			korean++
		case isCyrillic(r):
			cyrillic++
		case isArabic(r):
			arabic++
		case isHebrew(r):
			hebrew++
		}
	}

	// Determine primary language based on character counts
	total := len([]rune(text))
	if total == 0 {
		return ld.defaultLanguage
	}

	// Check for specific language indicators
	if japanese > 0 {
		return "ja"
	}
	if korean > 0 {
		return "ko"
	}
	if chinese > total/3 { // Chinese needs more characters to be confident
		return "zh"
	}
	if cyrillic > latin {
		return "ru"
	}
	if arabic > latin {
		return "ar"
	}
	if hebrew > latin {
		return "he"
	}

	// Check for common language patterns
	lowerText := strings.ToLower(text)
	
	// French indicators
	if containsAny(lowerText, []string{" le ", " la ", " les ", " de ", " du ", " des ", "ç", "à", "è", "é", "ê"}) {
		return "fr"
	}
	
	// Spanish indicators
	if containsAny(lowerText, []string{" el ", " la ", " los ", " las ", " de ", " del ", "ñ", "¿", "¡"}) {
		return "es"
	}
	
	// German indicators
	if containsAny(lowerText, []string{" der ", " die ", " das ", " den ", " dem ", "ä", "ö", "ü", "ß"}) {
		return "de"
	}
	
	// Italian indicators
	if containsAny(lowerText, []string{" il ", " la ", " lo ", " gli ", " le ", " di ", " del ", " della "}) {
		return "it"
	}
	
	// Portuguese indicators
	if containsAny(lowerText, []string{" o ", " a ", " os ", " as ", " de ", " do ", " da ", "ã", "õ", "ç"}) {
		return "pt"
	}

	// Default to English for Latin script
	if latin > total/2 {
		return "en"
	}

	return ld.defaultLanguage
}

// Character type detection functions
func isLatin(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		unicode.In(r, unicode.Latin)
}

func isJapanese(r rune) bool {
	return (r >= 0x3040 && r <= 0x309F) || // Hiragana
		(r >= 0x30A0 && r <= 0x30FF) || // Katakana
		(r >= 0xFF00 && r <= 0xFFEF) // Full-width
}

func isChinese(r rune) bool {
	return r >= 0x4E00 && r <= 0x9FFF
}

func isKorean(r rune) bool {
	return (r >= 0xAC00 && r <= 0xD7AF) || // Hangul Syllables
		(r >= 0x1100 && r <= 0x11FF) || // Hangul Jamo
		(r >= 0x3130 && r <= 0x318F) // Hangul Compatibility Jamo
}

func isCyrillic(r rune) bool {
	return unicode.In(r, unicode.Cyrillic)
}

func isArabic(r rune) bool {
	return unicode.In(r, unicode.Arabic)
}

func isHebrew(r rune) bool {
	return unicode.In(r, unicode.Hebrew)
}

// containsAny checks if the text contains any of the given substrings
func containsAny(text string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(text, substr) {
			return true
		}
	}
	return false
}