package lyrics

import "testing"

func TestCleanArtist(t *testing.T) {
	tests := []struct{ input, want string }{
		{"Dominic Fike", "Dominic Fike"},
		{"Dominic Fike - Topic", "Dominic Fike"},
		{"The Weeknd VEVO", "The Weeknd"},
		{"F.K.A. Twigs", "FKA Twigs"},
		{"Bad Bunny Ft. Someone", "Bad Bunny"},
		{"Bad Bunny feat. Someone", "Bad Bunny"},
	}
	for _, tt := range tests {
		if got := CleanArtist(tt.input); got != tt.want {
			t.Errorf("CleanArtist(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestCleanTrack(t *testing.T) {
	tests := []struct{ input, want string }{
		{"Babydoll", "Babydoll"},
		{"Cellophane (Official Video)", "Cellophane"},
		{"EARFQUAKE (Official Audio)", "EARFQUAKE"},
		{"HUMBLE.", "HUMBLE"},
	}
	for _, tt := range tests {
		if got := CleanTrack(tt.input); got != tt.want {
			t.Errorf("CleanTrack(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
