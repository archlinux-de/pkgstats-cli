package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pkgstats-cli/internal/build"
	"sort"
	"strconv"
	"time"
)

type Client struct {
	Client  *http.Client
	baseURL string
}

type PackagePopularity struct {
	Name       string
	Popularity float64
}

type PackagePopularityList struct {
	Total               int
	Count               int
	PackagePopularities []PackagePopularity
}

type packagePopularityResult struct {
	pp  PackagePopularity
	err error
}

func NewClient(baseURL string) *Client {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	client := Client{}
	client.Client = httpClient
	client.baseURL = baseURL

	return &client
}

func (client *Client) query(path string, params url.Values) ([]byte, error) {
	u, err := url.Parse(client.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("pkgstats/%s", build.Version))
	response, err := client.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	return body, err
}

func (client *Client) GetPackages(packages ...string) (PackagePopularityList, error) {
	ppl := PackagePopularityList{}
	ppl.PackagePopularities = make([]PackagePopularity, len(packages))

	ch := make(chan (packagePopularityResult))

	for _, p := range packages {
		go func(p string, ch chan (packagePopularityResult)) {
			res, err := client.GetPackage(p)
			ch <- packagePopularityResult{res, err}
		}(p, ch)
	}

	for i := range packages {
		res := <-ch
		if res.err != nil {
			return ppl, res.err
		}
		ppl.PackagePopularities[i] = res.pp
	}

	sort.Slice(ppl.PackagePopularities, func(i, j int) bool {
		return ppl.PackagePopularities[i].Popularity > ppl.PackagePopularities[j].Popularity
	})

	return ppl, nil
}

func (client *Client) GetPackage(p string) (PackagePopularity, error) {
	var pp PackagePopularity

	response, err := client.query(fmt.Sprintf("/api/packages/%s", url.QueryEscape(p)), url.Values{})

	if err != nil {
		return pp, err
	}

	err = json.Unmarshal(response, &pp)
	if err != nil {
		return pp, err
	}

	return pp, nil
}

func (client *Client) SearchPackages(query string, limit int) (PackagePopularityList, error) {
	var ppl PackagePopularityList

	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("query", query)

	response, err := client.query("/api/packages", params)

	if err != nil {
		return ppl, err
	}

	err = json.Unmarshal(response, &ppl)
	if err != nil {
		return ppl, err
	}

	return ppl, err
}
