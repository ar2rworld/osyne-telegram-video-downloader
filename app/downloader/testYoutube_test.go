package downloader

import (
	"testing"
)

func TestCutString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{"Empty string", "", 5, ""},
		{"String shorter than max", "Hello", 10, "Hello"},
		{"String equal to max", "Hello", 5, "Hello"},
		{"String longer than max", "Hello, World!", 5, "Hello"},
		{"Max length 0", "Hello", 0, ""},
		{"Unicode characters", "你好世界", 2, "\xe4\xbd"}, //nolint:gosmopolitan
		{"Max length 1", "Hello", 1, "H"},
		{"Spaces in string", "Hello World", 6, "Hello "},
		{"Very long string", "This is a very long string", 10, "This is a "},
		{"Negative max length", "Hello", -5, "Hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cutString(tt.input, tt.maxLength)
			if result != tt.expected {
				t.Errorf("cutString(%q, %d) = %q; want %q", tt.input, tt.maxLength, result, tt.expected)
			}
		})
	}
}

func TestConvertToUTF8(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"ASCII only", "Hello, World!", "Hello, World!"},
		{"UTF-8 characters", "こんにちは世界", "こんにちは世界"}, //nolint:gosmopolitan
		{"Mixed valid and invalid", "Hello\xFFWorld", "HelloWorld"},
		{"All invalid", "\xF0\x8C\xBC", ""},
		{"Empty string", "", ""},
		{"Emoji", "👋🌍", "👋🌍"},
		{"Mixed script", "Hello världen こんにちは 123", "Hello världen こんにちは 123"},
		{"Invalid sequence with valid ASCII", "\xf0(\x8c\xbc", "("},
		{"Multiple invalid sequences", "a\xFFb\xFEc\xFD", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToUTF8(tt.input)
			if result != tt.expected {
				t.Errorf("%s ConvertToUTF8(%q) = %q; want %q", tt.name, tt.input, result, tt.expected)
			}
		})
	}
}
