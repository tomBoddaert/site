package main

import (
	"unicode/utf8"
)

func truncate(s string, maxLen uint) string {
	// If the length in bytes of s < maxLen, then
	// the length in runes of s < maxLen
	if len(s) <= int(maxLen) {
		return s
	}

	if maxLen == 0 {
		return ""
	}

	var rlen uint = 0
	for i, r := range s {
		rlen++

		if rlen == maxLen {
			truncated := s[:i+utf8.RuneLen(r)]

			// Terminate the string with "..."
			terminated := []byte(truncated)
			terminated = append(terminated, '.', '.', '.')

			return string(terminated)
		}
	}

	// length in runes of s < maxLen
	return s
}
