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
		debugLog("Auto-detect disabled, using default language: %s", ld.defaultLanguage)
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
		debugLog("Empty text, using default language: %s", ld.defaultLanguage)
		return ld.defaultLanguage
	}

	// Check for specific language indicators
	if japanese > 0 {
		confidence := float64(japanese) / float64(total)
		debugLogLanguageDetection(text, "ja", confidence)
		return "ja"
	}
	if korean > 0 {
		confidence := float64(korean) / float64(total)
		debugLogLanguageDetection(text, "ko", confidence)
		return "ko"
	}
	if chinese > total/3 { // Chinese needs more characters to be confident
		confidence := float64(chinese) / float64(total)
		debugLogLanguageDetection(text, "zh", confidence)
		return "zh"
	}
	if cyrillic > latin {
		confidence := float64(cyrillic) / float64(total)
		debugLogLanguageDetection(text, "ru", confidence)
		return "ru"
	}
	if arabic > latin {
		confidence := float64(arabic) / float64(total)
		debugLogLanguageDetection(text, "ar", confidence)
		return "ar"
	}
	if hebrew > latin {
		confidence := float64(hebrew) / float64(total)
		debugLogLanguageDetection(text, "he", confidence)
		return "he"
	}

	// Check for common language patterns
	lowerText := strings.ToLower(text)
	
	// Check for more specific indicators first
	// Spanish indicators - check special characters first
	if containsAny(lowerText, []string{"ñ", "¿", "¡", "á", "é", "í", "ó", "ú"}) && containsAny(lowerText, []string{" el ", " la ", " los ", " las "}) {
		debugLogLanguageDetection(text, "es", 0.8)
		return "es"
	}
	if containsAny(lowerText, []string{"ñ", "¿", "¡"}) || 
		(containsAny(lowerText, []string{" el ", " los ", " las ", " del "}) && !containsAny(lowerText, []string{" le ", " les ", " du ", " des "})) {
		debugLogLanguageDetection(text, "es", 0.8)
		return "es"
	}
	
	// Portuguese indicators - check special characters first
	if containsAny(lowerText, []string{"ã", "õ", "ção"}) || 
		(containsAny(lowerText, []string{" o ", " os ", " as ", " do ", " da ", " na "}) && 
			containsAny(lowerText, []string{"á", "é", "ê", "ó", "ô"}) && 
			!containsAny(lowerText, []string{" el ", " la "})) {
		debugLogLanguageDetection(text, "pt", 0.8)
		return "pt"
	}
	
	// German indicators
	if containsAny(lowerText, []string{"ä", "ö", "ü", "ß"}) || 
		containsAny(lowerText, []string{" der ", " die ", " das ", " den ", " dem "}) {
		debugLogLanguageDetection(text, "de", 0.8)
		return "de"
	}
	
	// Italian indicators - more specific patterns
	if containsAny(lowerText, []string{" il ", " lo ", " gli ", " della ", " nel ", " è "}) {
		debugLogLanguageDetection(text, "it", 0.8)
		return "it"
	}
	
	// French indicators - check last to avoid conflicts
	if containsAny(lowerText, []string{"ç", "à", "è", "é", "ê", "ù"}) || 
		(containsAny(lowerText, []string{" le ", " les ", " du ", " des ", " sur ", " est "}) && !containsAny(lowerText, []string{" el ", " está "})) {
		debugLogLanguageDetection(text, "fr", 0.8)
		return "fr"
	}

	// Default to English for Latin script
	if latin > total/2 {
		confidence := float64(latin) / float64(total)
		debugLogLanguageDetection(text, "en", confidence)
		return "en"
	}

	debugLogLanguageDetection(text, ld.defaultLanguage, 0.0)
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