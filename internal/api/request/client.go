package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"pkgstats-cli/internal/build"
)

const timeout = 5 * time.Second

type Client struct {
	Client  *http.Client
	BaseURL string
}

type PackagePopularity struct {
	Name       string  `json:"name"`
	Popularity float64 `json:"popularity"`
}

type PackagePopularityList struct {
	Total               int                 `json:"total"`
	Count               int                 `json:"count"`
	PackagePopularities []PackagePopularity `json:"packagePopularities"`
}

type packagePopularityResult struct {
	pp  PackagePopularity
	err error
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

func (client *Client) query(path string, params url.Values) ([]byte, error) {
	u, err := url.Parse(client.BaseURL)
	if err != nil {
		return nil, err
	}
	u = u.JoinPath(path)

	u.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", build.UserAgent)
	response, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(response.StatusCode))
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
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
			return PackagePopularityList{}, res.err
		}
		ppl.PackagePopularities[i] = res.pp
	}

	sort.Slice(ppl.PackagePopularities, func(i, j int) bool {
		return ppl.PackagePopularities[i].Popularity > ppl.PackagePopularities[j].Popularity
	})

	return ppl, nil
}

func (client *Client) GetPackage(p string) (PackagePopularity, error) {
	response, err := client.query("/api/packages/"+url.PathEscape(p), url.Values{})
	if err != nil {
		return PackagePopularity{}, err
	}

	var pp PackagePopularity
	err = json.Unmarshal(response, &pp)
	if err != nil {
		return PackagePopularity{}, err
	}

	return pp, nil
}

func (client *Client) SearchPackages(query string, limit int) (PackagePopularityList, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("query", query)

	response, err := client.query("/api/packages", params)
	if err != nil {
		return PackagePopularityList{}, err
	}

	var ppl PackagePopularityList
	err = json.Unmarshal(response, &ppl)
	if err != nil {
		return PackagePopularityList{}, err
	}

	return ppl, nil
}
