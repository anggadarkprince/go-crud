package utilities

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Trim and normalize spaces
func CleanString(s string) string {
    return strings.Join(strings.Fields(s), " ")
}

// Convert string to title case
var titleCaser = cases.Title(language.English)
var titleCaserID = cases.Title(language.Indonesian)
func TitleCase(s string) string {
    return titleCaser.String(strings.ToLower(s))
}
func TitleCaseID(s string) string {
    return titleCaserID.String(strings.ToLower(s))
}
func Capitalize(s string) string {
    if len(s) == 0 {
        return s
    }
    s = strings.ReplaceAll(s, "_", " ")
    s = strings.ReplaceAll(s, "-", " ")
    return strings.ToUpper(s[:1]) + s[1:]
}

// Generate secure random string
func RandomString(n int) string {
    b := make([]byte, n)
    _, _ = rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)[:n]
}

// Check if string in slice
func InSlice(str string, list []string) bool {
    for _, v := range list {
        if v == str {
            return true
        }
    }
    return false
}