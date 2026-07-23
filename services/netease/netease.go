// Package netease is a lyrics.Provider backed by NetEase Cloud Music's
// public search/lyric endpoints, used as a fallback when lrclib has nothing.
package netease

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"slyrics/lyrics"
)

const userAgent = "Mozilla/5.0 (compatible; Slyrics)"

func New() *Client {
	return &Client{}
}

type Client struct {
	http http.Client
}

// Client implements lyrics.Provider
func (c *Client) Lyrics(artist, track string) ([]lyrics.Line, error) {
	artist, track = lyrics.CleanArtist(artist), lyrics.CleanTrack(track)
	id, err := c.searchID(artist, track)
	if err != nil || id == 0 {
		return nil, err
	}
	return c.lyric(id)
}

// searchID looks up candidates by title and picks the first one whose

func (c *Client) searchID(artist, track string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	u := "https://music.163.com/api/search/get?" + url.Values{
		"s":     {strings.TrimSpace(artist + " " + track)},
		"type":  {"1"},
		"limit": {"10"},
	}.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var response searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	for _, song := range response.Result.Songs {
		if artist == "" {
			return song.ID, nil
		}
		for _, a := range song.Artists {
			if lyrics.SameArtist(artist, a.Name) {
				return song.ID, nil
			}
		}
	}
	return 0, nil
}

func (c *Client) lyric(id int) ([]lyrics.Line, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	u := "https://music.163.com/api/song/lyric?" + url.Values{
		"id": {strconv.Itoa(id)},
		"lv": {"1"},
		"kv": {"1"},
		"tv": {"-1"},
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

	var response lyricResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return parseLrc(response.Lrc.Lyric), nil
}

func parseLrc(raw string) []lyrics.Line {
	if raw == "" {
		return nil
	}
	lines := strings.Split(raw, "\n")
	result := make([]lyrics.Line, 0, len(lines))
	for _, line := range lines {
		if !lyrics.IsTimestampLine(line) {
			continue
		}
		result = append(result, lyrics.ParseLrcLine(line))
	}
	return result
}

type searchResponse struct {
	Result struct {
		Songs []struct {
			ID      int `json:"id"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"songs"`
	} `json:"result"`
}

type lyricResponse struct {
	Lrc struct {
		Lyric string `json:"lyric"`
	} `json:"lrc"`
}
