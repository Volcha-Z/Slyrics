package lrclib

import "testing"

func TestSearchDoesNotReturnWrongArtist(t *testing.T) {
	lines, err := New().Lyrics("Dominic Fike", "Babydoll")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) == 0 {
		t.Fatal("expected to find lyrics for a well-known track")
	}
	found := false
	for _, l := range lines {
		if l.Words == "I can't move on, babydoll" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("got lyrics that don't look like Dominic Fike's Babydoll: %+v", lines)
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
