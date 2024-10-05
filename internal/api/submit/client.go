package submit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"pkgstats-cli/internal/build"
)

const timeout = 15 * time.Second

type Client struct {
	Client  *http.Client
	BaseURL string
}

func NewClient(baseURL string) *Client {
	httpClient := &http.Client{
		Timeout: timeout,
	}

	return &Client{
		Client:  httpClient,
		BaseURL: baseURL,
	}
}

func (client *Client) SendRequest(request Request) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	u, err := url.Parse(client.BaseURL)
	if err != nil {
		return err
	}
	u = u.JoinPath("/api/submit")

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("pkgstats/%s", build.Version))
	response, err := client.Client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(response.Body)
		err = fmt.Errorf("server error: %s", string(body))
	}

	return err
}
