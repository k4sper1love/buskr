package mdutil

import "html"

// Escape escapes HTML special characters in user-provided strings
func Escape(s string) string {
	return html.EscapeString(s)
}
