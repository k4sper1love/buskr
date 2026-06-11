package mapimg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/k4sper1love/buskr/internal/domain/location"
)

const (
	staticMapsURL = "https://maps.googleapis.com/maps/api/staticmap"
	mapSize       = "640x640"
	mapScale      = "2"
	maxLabeled    = 26 // A–Z
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// Generate returns the PNG image bytes and a numbered list "A. Name\nB. Name\n..."
func (c *Client) Generate(locs []*location.Location) ([]byte, string, error) {
	if c.apiKey == "" {
		return nil, "", fmt.Errorf("google maps api key is not configured")
	}

	params := url.Values{}
	params.Set("size", mapSize)
	params.Set("scale", mapScale)
	params.Set("maptype", "roadmap")
	params.Set("key", c.apiKey)

	var sb strings.Builder
	for i, loc := range locs {
		var marker string
		if i < maxLabeled {
			label := string(rune('A' + i))
			marker = fmt.Sprintf("color:red|label:%s|%f,%f", label, loc.Coords.Lat, loc.Coords.Lon)
			fmt.Fprintf(&sb, "%s. %s\n", label, loc.Name)
		} else {
			marker = fmt.Sprintf("color:red|%f,%f", loc.Coords.Lat, loc.Coords.Lon)
		}
		params.Add("markers", marker)
	}

	resp, err := c.httpClient.Get(staticMapsURL + "?" + params.Encode())
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("google maps api: status %d", resp.StatusCode)
	}

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return imgBytes, strings.TrimRight(sb.String(), "\n"), nil
}
