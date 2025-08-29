package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	"pkgstats-cli/internal/build"
)

const (
	timeout        = 5 * time.Second
	maxConcurrency = 4
)

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

func (client *Client) GetPackages(packages ...string) (PackagePopularityList, []error) {
	var wg sync.WaitGroup
	var errs []error
	var mu sync.Mutex

	jobs := make(chan string, len(packages))
	results := make(chan packagePopularityResult, len(packages))

	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				res, err := client.GetPackage(p)
				results <- packagePopularityResult{res, err}
			}
		}()
	}

	for _, p := range packages {
		jobs <- p
	}
	close(jobs)

	wg.Wait()
	close(results)

	ppl := PackagePopularityList{}
	for res := range results {
		if res.err != nil {
			mu.Lock()
			errs = append(errs, res.err)
			mu.Unlock()
		} else {
			ppl.PackagePopularities = append(ppl.PackagePopularities, res.pp)
		}
	}

	sort.Slice(ppl.PackagePopularities, func(i, j int) bool {
		return ppl.PackagePopularities[i].Popularity > ppl.PackagePopularities[j].Popularity
	})

	return ppl, errs
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
