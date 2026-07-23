package lrclib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"slyrics/lyrics"
)

const userAgent = "Slyrics/1.0"

func New() *Client {
	return &Client{}
}

type Client struct {
	http http.Client
}

// Client implements lyrics.Provider
func (c *Client) Lyrics(artist, track string) ([]lyrics.Line, error) {
	artist, track = lyrics.CleanArtist(artist), lyrics.CleanTrack(track)
	if artist == "" || track == "" {
		return c.search(artist, track)
	}

	// the exact-match endpoint is strict about metadata (feat. artists,
	// remaster suffixes, etc.) and misses tracks the fuzzy search finds,
	// so fall back to search instead of reporting no lyrics.
	lines, err := c.get(artist, track)
	if err != nil {
		return nil, err
	}
	if len(lines) > 0 {
		return lines, nil
	}
	return c.search(artist, track)
}

func (c *Client) get(artist, track string) ([]lyrics.Line, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u := "https://lrclib.net/api/get?" + url.Values{
		"artist_name": {artist},
		"track_name":  {track},
	}.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var response lrclibTrack
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return parseTrack(response), nil
}

// search hits lrclib's fuzzy-text endpoint, which ranks by title similarity
// and ignores artist entirely — for a common title the top hits are often
// other artists' songs. If we know the artist, only a matching result is
// trusted; otherwise we'd rather report no lyrics than the wrong song's.
func (c *Client) search(artist, track string) ([]lyrics.Line, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u := "https://lrclib.net/api/search?" + url.Values{
		"q": {strings.TrimSpace(artist + " " + track)},
	}.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response []lrclibTrack
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	for _, t := range response {
		if artist != "" && !lyrics.SameArtist(artist, t.ArtistName) {
			continue
		}
		if lines := parseTrack(t); len(lines) > 0 {
			return lines, nil
		}
	}
	return nil, nil
}

type lrclibTrack struct {
	ArtistName   string `json:"artistName"`
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"`
}

func parseTrack(t lrclibTrack) []lyrics.Line {
	if t.SyncedLyrics != "" {
		return parseSynced(t)
	}
	if t.PlainLyrics != "" {
		return parsePlain(t)
	}
	return nil
}

func parseSynced(r lrclibTrack) []lyrics.Line {
	lines := strings.Split(r.SyncedLyrics, "\n")
	result := make([]lyrics.Line, 0, len(lines))
	for _, line := range lines {
		if !lyrics.IsTimestampLine(line) {
			continue
		}
		result = append(result, lyrics.ParseLrcLine(line))
	}
	return result
}

func parsePlain(r lrclibTrack) []lyrics.Line {
	lines := strings.Split(r.PlainLyrics, "\n")
	result := make([]lyrics.Line, len(lines))
	for i, line := range lines {
		result[i] = lyrics.Line{Words: line}
	}
	return result
}
