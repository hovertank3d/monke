package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hovertank3d/monke"
	"github.com/hovertank3d/monke/env"
)

var host = env.MonkeHost + "/api"
var vidHost = env.MonkeHost + "/proxy"

var DefaultClient = &Client{}

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) getJSON(dst interface{}, path string, params url.Values) error {
	resp, err := http.DefaultClient.Get(host + path + "?" + params.Encode())
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(dst)
}

func (c *Client) Search(tags string, name string) (animes []monke.Anime, _ error) {
	p := url.Values{}
	p.Set("tags", tags)
	p.Set("name", name)
	return animes, c.getJSON(&animes, "/search", p)
}

func (c *Client) Anime(id string) (anime monke.Anime, _ error) {
	return anime, c.getJSON(&anime, "/"+id, nil)
}

func (c *Client) VideoLink(id string, episode int, season int) string {
	return fmt.Sprintf("%s/%s/%d/%d.mp4", vidHost, id, episode, season)
}
