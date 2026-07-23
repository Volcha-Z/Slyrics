package netease

import "testing"

func TestLyrics(t *testing.T) {
	lines, err := New().Lyrics("The Weeknd", "Blinding Lights")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) == 0 {
		t.Fatal("expected non-empty lyrics for a well-known track")
	}
}

func TestSearchRejectsWrongArtist(t *testing.T) {
	lines, err := New().Lyrics("A Completely Made Up Artist Name", "Babydoll")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 {
		t.Errorf("expected no match for a nonexistent artist, got %d lines", len(lines))
	}
}
