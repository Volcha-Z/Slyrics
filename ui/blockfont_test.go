package ui

import (
	"strings"
	"testing"
)

func TestRenderBlockTextWrapsWithinWidth(t *testing.T) {
	const maxWidth = 60
	out := renderBlockText("I had it all along I had it all along I had it", maxWidth, true)
	lines := strings.Split(out, "\n")

	groupHeight := blockFontHeight + 1
	if (len(lines)+1)%(groupHeight+1) != 0 {
		t.Fatalf("unexpected line count for %d-row groups with spacers, got %d", groupHeight, len(lines))
	}
	for i, l := range lines {
		isSpacer := (i+1)%(groupHeight+1) == 0
		if isSpacer != (l == "") {
			t.Errorf("line %d: expected spacer=%v, got %q", i, isSpacer, l)
		}
		if got := len([]rune(l)); got > maxWidth {
			t.Errorf("line %d exceeds maxWidth %d: %q (%d cols)", i, maxWidth, l, got)
		}
	}
}

func TestRenderBlockTextEmpty(t *testing.T) {
	if out := renderBlockText("", 40, true); out != "" {
		t.Errorf("expected empty output for empty input, got %q", out)
	}
}

func TestWidestWordGlyphWidth(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"single word", "HELLO"},
		{"picks the longer word", "HI THOUGHT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := widestWordGlyphWidth(tt.input)
			words := strings.Fields(tt.input)
			longest := 0
			for _, w := range words {
				if l := len(w); l > longest {
					longest = l
				}
			}
			if got <= 0 || got < longest {
				t.Errorf("widestWordGlyphWidth(%q) = %d, want at least %d", tt.input, got, longest)
			}
		})
	}
}
