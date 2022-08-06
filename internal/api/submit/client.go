package submit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pkgstats-cli/internal/build"
	"time"
)

type Client struct {
	Client  *http.Client
	baseURL string
}

func NewClient(baseURL string) *Client {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	client := Client{}
	client.Client = httpClient
	client.baseURL = baseURL

	return &client
}

func (client *Client) SendRequest(request Request) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	u, err := url.Parse(client.baseURL)
	if err != nil {
		return err
	}
	u.Path = "/api/submit"

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(payload))
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

	if response.StatusCode != 204 && err == nil {
		body, _ := io.ReadAll(response.Body)
		err = errors.New("Server Error:" + string(body))
	}

	defer response.Body.Close()

	return err
}
