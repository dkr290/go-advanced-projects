package chunk

import "strings"

func chunk(text string, size, overlap int) []string {
	text = strings.TrimSpace(text)
	if text == "" || size <= 0 {
		return nil
	}

	runes := []rune(text)
	if len(runes) <= size {
		return []string{text}
	}

	if overlap < 0 {
		overlap = 0
	}
	if overlap >= size {
		overlap = size / 2
	}
	step := size - overlap

	var chunks []string
	i := 0
	for i < len(runes) {
		end := min(i+size, len(runes))

		// Try to find a boundary (space/newline/tab) so we don't cut mid-word
		if end < len(runes) && !isBoundary(runes[end]) {
			if b := lastIndexBoundary(runes, i, end); b > i {
				end = b
			}
		}

		if part := strings.TrimSpace(string(runes[i:end])); part != "" {
			chunks = append(chunks, part)
		}

		if end >= len(runes) {
			break
		}

		// Advance by step, but never less than what we actually emitted,
		// so overlap can't push us backward or stall the loop.
		advance := step
		if emitted := end - i; emitted < advance {
			advance = emitted
		}
		if advance <= 0 {
			advance = 1
		}
		i += advance
	}

	return chunks
}

func isBoundary(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

func lastIndexBoundary(runes []rune, start, end int) int {
	for j := end - 1; j > start; j-- {
		if isBoundary(runes[j]) {
			return j
		}
	}
	return -1
}
