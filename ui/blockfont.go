package ui

import (
	"strings"

	gloss "github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

var blockLetters = map[rune][]string{
	'A':  {"  ███  ", " ██ ██ ", "███████", "██   ██", "██   ██"},
	'B':  {"██████ ", "██   ██", "██████ ", "██   ██", "██████ "},
	'C':  {" █████ ", "██   ██", "██     ", "██   ██", " █████ "},
	'D':  {"██████ ", "██   ██", "██   ██", "██   ██", "██████ "},
	'E':  {"███████", "██     ", "█████  ", "██     ", "███████"},
	'F':  {"███████", "██     ", "█████  ", "██     ", "██     "},
	'G':  {" █████ ", "██     ", "██  ███", "██   ██", " █████ "},
	'H':  {"██   ██", "██   ██", "███████", "██   ██", "██   ██"},
	'I':  {"████", " ██ ", " ██ ", " ██ ", "████"},
	'J':  {"     ██", "     ██", "     ██", "██   ██", " █████ "},
	'K':  {"██   ██", "██  ██ ", "█████  ", "██  ██ ", "██   ██"},
	'L':  {"██     ", "██     ", "██     ", "██     ", "███████"},
	'M':  {"██   ██", "███ ███", "███████", "██ █ ██", "██   ██"},
	'N':  {"██   ██", "███  ██", "████ ██", "██ ████", "██   ██"},
	'O':  {" █████ ", "██   ██", "██   ██", "██   ██", " █████ "},
	'P':  {"██████ ", "██   ██", "██████ ", "██     ", "██     "},
	'Q':  {" █████ ", "██   ██", "██   ██", "██  ███", " ██████"},
	'R':  {"██████ ", "██   ██", "██████ ", "██  ██ ", "██   ██"},
	'S':  {" █████ ", "██     ", " █████ ", "     ██", " █████ "},
	'T':  {"███████", "  ██   ", "  ██   ", "  ██   ", "  ██   "},
	'U':  {"██   ██", "██   ██", "██   ██", "██   ██", " █████ "},
	'V':  {"██   ██", "██   ██", "██   ██", " ██ ██ ", "  ███  "},
	'W':  {"██   ██", "██   ██", "██ █ ██", "███████", "███ ███"},
	'X':  {"██   ██", " ██ ██ ", "  ███  ", " ██ ██ ", "██   ██"},
	'Y':  {"██   ██", " ██ ██ ", "  ███  ", "  ██   ", "  ██   "},
	'Z':  {"███████", "    ██ ", "  ███  ", " ██    ", "███████"},
	' ':  {"    ", "    ", "    ", "    ", "    "},
	'0':  {" █████ ", "██   ██", "██   ██", "██   ██", " █████ "},
	'1':  {"  ██   ", " ███   ", "  ██   ", "  ██   ", "███████"},
	'2':  {" █████ ", "██   ██", "   ███ ", " ██    ", "███████"},
	'3':  {" █████ ", "██   ██", "  ████ ", "██   ██", " █████ "},
	'4':  {"██   ██", "██   ██", "███████", "     ██", "     ██"},
	'5':  {"███████", "██     ", "██████ ", "     ██", "██████ "},
	'6':  {" █████ ", "██     ", "██████ ", "██   ██", " █████ "},
	'7':  {"███████", "     ██", "    ██ ", "   ██  ", "  ██   "},
	'8':  {" █████ ", "██   ██", " █████ ", "██   ██", " █████ "},
	'9':  {" █████ ", "██   ██", " ██████", "     ██", " █████ "},
	'\'': {"██", "██", "  ", "  ", "  "},
	',':  {"  ", "  ", "  ", "██", "█ "},
	'.':  {"  ", "  ", "  ", "  ", "██"},
	'!':  {"██", "██", "██", "  ", "██"},
	'?':  {" ███ ", "█   █", "   █ ", "     ", "  █  "},
	'-':  {"      ", "      ", "██████", "      ", "      "},
	'(':  {" ██", "██ ", "██ ", "██ ", " ██"},
	')':  {"██ ", " ██", " ██", " ██", "██ "},
	':':  {"  ", "██", "  ", "██", "  "},
	';':  {"  ", "██", "  ", "██", "█ "},
}

const blockFontHeight = 5

// gap between words, wider than the gap between letters
const wordGap = 3

func blockGlyph(c rune) []string {
	if g, ok := blockLetters[c]; ok {
		return g
	}
	return blockLetters[' ']
}

func glyphWidth(c rune) int {
	return len([]rune(blockGlyph(c)[0])) + 1
}

// widestWordGlyphWidth is the width a caller needs before block-font
// rendering is even worth trying: renderBlockText only wraps at word
// boundaries, so a word wider than the window overflows it.
func widestWordGlyphWidth(s string) int {
	max := 0
	for _, word := range strings.Fields(strings.ToUpper(s)) {
		w := 0
		for _, c := range word {
			w += glyphWidth(c)
		}
		if w > max {
			max = w
		}
	}
	return max
}

func renderBlockText(s string, maxWidth int, effect3D bool) string {
	words := strings.Fields(strings.ToUpper(s))
	if len(words) == 0 || maxWidth < 1 {
		return ""
	}
	if effect3D {
		maxWidth--
	}

	var rows [][]string
	current := make([]string, blockFontHeight)
	currentWidth := 0

	flush := func() {
		if currentWidth == 0 {
			return
		}
		if effect3D {
			rows = append(rows, shadow3D(current))
		} else {
			rows = append(rows, current)
		}
		current = make([]string, blockFontHeight)
		currentWidth = 0
	}

	for _, word := range words {
		wordLines := make([]string, blockFontHeight)
		wordWidth := 0
		for _, c := range word {
			g := blockGlyph(c)
			for i := 0; i < blockFontHeight; i++ {
				wordLines[i] += g[i] + " "
			}
			wordWidth += glyphWidth(c)
		}

		gap := 0
		if currentWidth > 0 {
			gap = wordGap
		}
		if currentWidth > 0 && currentWidth+gap+wordWidth > maxWidth {
			flush()
			gap = 0
		}

		if currentWidth > 0 {
			for i := 0; i < blockFontHeight; i++ {
				current[i] += strings.Repeat(" ", wordGap) + wordLines[i]
			}
			currentWidth += gap + wordWidth
		} else {
			current = wordLines
			currentWidth = wordWidth
		}
	}
	flush()

	var lines []string
	for i, row := range rows {
		if i > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, row...)
	}
	return strings.Join(lines, "\n")
}

// 3D drop shadow behind block glyphs
func shadow3D(lines []string) []string {
	h := len(lines) + 1
	w := 0
	for _, l := range lines {
		if rl := len([]rune(l)); rl > w {
			w = rl
		}
	}
	w++

	grid := make([][]rune, h)
	for i := range grid {
		grid[i] = make([]rune, w)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}
	for y, line := range lines {
		for x, c := range []rune(line) {
			if c == '█' {
				grid[y+1][x+1] = '░'
			}
		}
	}
	for y, line := range lines {
		for x, c := range []rune(line) {
			if c == '█' {
				grid[y][x] = '█'
			}
		}
	}

	out := make([]string, h)
	for i, row := range grid {
		out[i] = string(row)
	}
	return out
}

// alignBlockText pads each line to sit at align within width. Doesn't go
// through lipgloss's Width(): it measures block glyphs differently than we
// do here, and re-wrapping already-wrapped text splits rows mid-glyph.
func alignBlockText(s string, width int, align gloss.Position) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		pad := width - runewidth.StringWidth(line)
		if pad <= 0 {
			continue
		}
		left := int(float64(pad) * float64(align))
		lines[i] = strings.Repeat(" ", left) + line
	}
	return strings.Join(lines, "\n")
}
