package lyrics

import "strings"

// SameArtist does a loose comparison for matching search results against
// the artist we actually asked for. Search endpoints return whatever's
// textually close to the query, which for a generic title can easily be a
// different artist entirely — better to report no lyrics than confidently
// show the wrong song's.
func SameArtist(expected, candidate string) bool {
	expected = strings.ToLower(strings.TrimSpace(expected))
	candidate = strings.ToLower(strings.TrimSpace(candidate))
	if expected == "" || candidate == "" {
		return false
	}
	return strings.Contains(candidate, expected) || strings.Contains(expected, candidate)
}
