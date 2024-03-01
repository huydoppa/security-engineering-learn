package util

import (
	"os"
	"strings"
	"unicode"
)

var black_list = []string{
	"com\\android\\", "\\androidx\\", "\\org\\", "intellij",
	"chromium", "jetbrains", "google", "facebook", "org\\slf4j", "\\kotlin\\", "\\okhttp3\\",
	"\\android\\", "\\firebase\\", "\\javax\\", "\\resources\\"}

func Contains(s []string, e string) bool {
	for _, a := range s {

		if strings.Contains(e, a) {

			return true
		}
	}
	return false
}

func SkipFolder(folder string) bool {
	if Contains(black_list, folder) {
		return true
	}
	return false

}
func ContainsNonASCII(s string) bool {
	for _, char := range s {
		if char > 127 {
			return true
		}
	}
	return false
}
func HasUpper(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func IsReactNative(dir string) bool {
	if _, err := os.Stat(dir + "\\assets\\index.android.bundle"); err == nil {
		return true
	}
	return false
}
