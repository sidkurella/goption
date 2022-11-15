package stringutil

// Truncates a string by code point length.
// If you want to truncate by raw byte length, use TruncateBytes.
func Truncate(s string, maxLength int) string {
	r := []rune(s)
	if len(r) < maxLength {
		return s
	}
	return string(r[:maxLength])
}

// Truncates a string by byte length.
// If you want to truncate by code point length, use Truncate.
func TruncateBytes(s string, maxLength int) string {
	b := []byte(s)
	if len(b) < maxLength {
		return s
	}
	return string(b[:maxLength])
}
