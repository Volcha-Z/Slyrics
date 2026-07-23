//go:build windows || darwin

package mpris

import (
	"errors"
	"slyrics/player"
)

func New(players []string) (*Client, error) {
	return nil, errors.New("darwin is not supported")
}

// Client implements player.Player
type Client struct{}

func (p *Client) State() (*player.State, error) {
	return nil, nil
}
