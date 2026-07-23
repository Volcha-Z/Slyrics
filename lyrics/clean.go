package lyrics

import (
	"regexp"
	"strings"
)

// trackNoise strips the "(Official Video)"-style decoration YouTube titles
// carry, which otherwise pollutes lyric search queries enough that even the
// right track stops showing up in the results.
var trackNoise = regexp.MustCompile(`(?i)[\(\[][^)\]]*\b(official\s*(music\s*)?video|official\s*audio|lyric(s)?\s*video|visualizer|explicit)\b[^)\]]*[\)\]]`)

// artistSuffix strips YouTube channel conventions ("Artist - Topic",
// "Artist VEVO") that never appear in how lyric sites credit the artist.
var artistSuffix = regexp.MustCompile(`(?i)\s*-\s*topic\s*$|\s+vevo\s*$`)

// artistFeature cuts a query down to the primary artist — search engines
// rank worse, not better, when the query includes every featured artist.
var artistFeature = regexp.MustCompile(`(?i)\s+(ft\.?|feat\.?|featuring)\s+.*$`)

func CleanTrack(s string) string {
	s = trackNoise.ReplaceAllString(s, "")
	return collapseSpace(stripPeriods(s))
}

func CleanArtist(s string) string {
	s = artistSuffix.ReplaceAllString(s, "")
	s = artistFeature.ReplaceAllString(s, "")
	return collapseSpace(stripPeriods(s))
}

// stripPeriods drops periods from stylized names like "F.K.A. Twigs" —
// lrclib's search returns zero results outright for a query containing them.
func stripPeriods(s string) string {
	return strings.ReplaceAll(s, ".", "")
}

func collapseSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
