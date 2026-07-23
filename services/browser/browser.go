package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slyrics/player"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// debug prints every raw message the extension sends, for diagnosing
// connections that establish but never seem to report anything.
var debug = os.Getenv("SLYRICS_DEBUG") != ""

const helloMessage = "ADAPTER_VERSION 1.0.0;WNPRLIB_REVISION 2"

type state int

const (
	stopped state = iota
	paused
	playing
)

func New(port int) (*Client, error) {
	c := &Client{}
	return c, c.start(port)
}

// Client implements player.Player
type Client struct {
	state    state
	position int
	title    string
	artist   string
	coverURL string

	updateTime time.Time

	stateMu sync.Mutex
	connMu  sync.Mutex

	// artist resolution runs in the background and is cached per video, so
	// State() never blocks on the network round trip
	resolveMu      sync.Mutex
	resolvedFor    string
	resolvedArtist string
	resolvingFor   string
}

func (c *Client) handler(w http.ResponseWriter, r *http.Request) {
	// make sure we only have one connection
	c.connMu.Lock()
	defer c.connMu.Unlock()

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		if debug {
			fmt.Fprintln(os.Stderr, "[browser] accept failed:", err)
		}
		return
	}
	defer conn.Close(websocket.StatusInternalError, "internal error")

	if debug {
		fmt.Fprintln(os.Stderr, "[browser] client connected from", r.RemoteAddr)
	}

	writer, err := conn.Writer(r.Context(), websocket.MessageText)
	if err != nil {
		return
	}

	writer.Write([]byte(helloMessage))
	writer.Close()

	for {
		t, reader, err := conn.Reader(r.Context())
		if err != nil {
			if debug {
				fmt.Fprintln(os.Stderr, "[browser] connection closed:", err)
			}
			return
		}

		msg, err := io.ReadAll(reader)
		if err != nil {
			return
		}
		if t != websocket.MessageText || len(msg) == 0 {
			continue
		}
		if debug {
			fmt.Fprintln(os.Stderr, "[browser] recv:", string(msg))
		}
		c.processMessage(string(msg))
	}
}

func (c *Client) processMessage(msg string) {
	spaceIndex := strings.IndexByte(msg, ' ')
	if spaceIndex == -1 {
		return
	}

	msgType := strings.ToUpper(msg[:spaceIndex])
	data := msg[spaceIndex+1:]

	// we are not doing global locking here because
	// we are not interested in most of the messages
	switch msgType {
	case "STATE":
		c.stateMu.Lock()
		switch data {
		case "PLAYING":
			c.state = playing
		case "PAUSED":
			c.state = paused
		case "STOPPED":
			c.state = stopped
		}
		c.stateMu.Unlock()
	case "TITLE":
		c.stateMu.Lock()
		c.title = data
		c.stateMu.Unlock()
	case "ARTIST":
		c.stateMu.Lock()
		c.artist = data
		c.stateMu.Unlock()
	case "COVER_URL":
		c.stateMu.Lock()
		c.coverURL = data
		c.stateMu.Unlock()
	case "POSITION_SECONDS":
		pos, _ := strconv.Atoi(data)
		c.stateMu.Lock()
		c.position = pos * 1000
		c.updateTime = time.Now()
		c.stateMu.Unlock()
	}
}

func (c *Client) start(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler: http.HandlerFunc(c.handler),
	}
	go server.Serve(l)
	return nil
}

func (c *Client) State() (*player.State, error) {
	c.stateMu.Lock()
	if c.state == stopped {
		c.stateMu.Unlock()
		return nil, nil
	}
	artist, title, coverURL := c.artist, c.title, c.coverURL
	position := c.position
	if c.state != paused {
		position += int(time.Since(c.updateTime).Milliseconds())
	}
	isPlaying := c.state == playing
	c.stateMu.Unlock()

	// YouTube localizes the artist name to the account's language (e.g.
	// Arabic script), which no lyrics source will have on file — the video
	// title itself is a better bet, since uploaders write it as plain
	// "Artist - Track" regardless of UI locale. Resolved in the background:
	// this runs on every poll, and blocking here would stall position
	// tracking by however long the request takes.
	if !hasLatinLetters(artist) {
		if resolved, ok := c.cachedArtist(coverURL); ok {
			artist = resolved
		} else {
			c.resolveArtistAsync(coverURL)
		}
	}

	return &player.State{
		ID:       artist + " " + title,
		Artist:   artist,
		Track:    title,
		Position: position,
		Playing:  isPlaying,
	}, nil
}

var videoIDPattern = regexp.MustCompile(`/vi/([^/]+)/`)

func (c *Client) cachedArtist(coverURL string) (string, bool) {
	c.resolveMu.Lock()
	defer c.resolveMu.Unlock()
	if c.resolvedFor == coverURL {
		return c.resolvedArtist, true
	}
	return "", false
}

func (c *Client) resolveArtistAsync(coverURL string) {
	c.resolveMu.Lock()
	if c.resolvingFor == coverURL {
		c.resolveMu.Unlock()
		return
	}
	c.resolvingFor = coverURL
	c.resolveMu.Unlock()

	go func() {
		var artist string
		if m := videoIDPattern.FindStringSubmatch(coverURL); m != nil {
			artist = fetchOEmbedArtist(m[1])
			if debug {
				fmt.Fprintf(os.Stderr, "[browser] resolved artist for video %s: %q\n", m[1], artist)
			}
		}
		c.resolveMu.Lock()
		c.resolvedFor = coverURL
		c.resolvedArtist = artist
		c.resolveMu.Unlock()
	}()
}

func fetchOEmbedArtist(videoID string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := "https://www.youtube.com/oembed?" + url.Values{
		"url":    {"https://www.youtube.com/watch?v=" + videoID},
		"format": {"json"},
	}.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var body struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return ""
	}

	if idx := strings.Index(body.Title, " - "); idx != -1 {
		return strings.TrimSpace(body.Title[:idx])
	}
	return ""
}

func hasLatinLetters(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}
