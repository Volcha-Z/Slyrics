package qqmusic

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
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

	song, err := c.search(artist, track)
	if err != nil || song == nil {
		return nil, err
	}
	return c.lyric(song.SongMID, song.SongID)
}

type songResult struct {
	SongMID string
	SongID  int
}

func (c *Client) search(artist, track string) (*songResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	u := "https://c.y.qq.com/soso/fcgi-bin/client_search_cp?" + url.Values{
		"w":      {strings.TrimSpace(artist + " " + track)},
		"format": {"json"},
		"p":      {"1"},
		"n":      {"10"},
	}.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://y.qq.com/")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	for _, song := range response.Data.Song.List {
		if artist == "" {
			return &songResult{song.SongMID, song.SongID}, nil
		}
		for _, singer := range song.Singer {
			if lyrics.SameArtist(artist, singer.Name) {
				return &songResult{song.SongMID, song.SongID}, nil
			}
		}
	}
	return nil, nil
}

func (c *Client) lyric(songMID string, songID int) ([]lyrics.Line, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	body, err := json.Marshal(map[string]any{
		"req_1": map[string]any{
			"method": "GetPlayLyricInfo",
			"module": "music.musichallSong.PlayLyricInfo",
			"param": map[string]any{
				"songMID": songMID,
				"songID":  songID,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://u.y.qq.com/cgi-bin/musicu.fcg", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://y.qq.com/")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response lyricResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	raw, err := base64.StdEncoding.DecodeString(response.Req1.Data.Lyric)
	if err != nil {
		return nil, err
	}
	return parseLrc(string(raw)), nil
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
	Data struct {
		Song struct {
			List []struct {
				SongMID string `json:"songmid"`
				SongID  int    `json:"songid"`
				Singer  []struct {
					Name string `json:"name"`
				} `json:"singer"`
			} `json:"list"`
		} `json:"song"`
	} `json:"data"`
}

type lyricResponse struct {
	Req1 struct {
		Data struct {
			Lyric string `json:"lyric"`
		} `json:"data"`
	} `json:"req_1"`
}
