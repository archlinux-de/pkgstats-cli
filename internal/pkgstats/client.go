package pkgstats

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"pkgstats-cli/internal/build"
	"strings"
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

	apiClient := Client{}
	apiClient.Client = httpClient
	apiClient.baseURL = baseURL

	return &apiClient
}

func (api *Client) SendRequest(packages string, cpuArchitecture string, architecture string, mirror string, quiet bool) (string, error) {
	form := url.Values{}
	form.Add("packages", packages)
	form.Add("arch", architecture)
	form.Add("cpuarch", cpuArchitecture)
	form.Add("mirror", mirror)
	if quiet {
		form.Add("quiet", "true")
	} else {
		form.Add("quiet", "false")
	}

	req, err := http.NewRequest("POST", api.baseURL+"/post", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("User-Agent", fmt.Sprintf("pkgstats/%s", build.Version))
	response, err := api.Client.Do(req)

	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 && err == nil {
		err = errors.New("Server Error")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	return strings.TrimSpace(string(body)), err
}
