package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	searchRepositoriesBaseURL = "https://api.github.com/search/repositories"
	maximumResults            = 5 // TODO: use 100
)

type Client struct {
	searchReposBaseURL string
	maxResults         uint
}

func NewClient() *Client {
	return &Client{
		searchReposBaseURL: searchRepositoriesBaseURL,
		maxResults:         maximumResults,
	}
}

func (c *Client) ListLatestRepositories(filters url.Values) ([]Repository, error) {
	searchURL := fmt.Sprintf(
		"%s?q=%s&sort=updated&order=desc&per_page=%d",
		c.searchReposBaseURL,
		intoSearchQuery(filters),
		c.maxResults,
	)

	body, err := httpRawGet(searchURL)
	if err != nil {
		return nil, err
	}

	var repos RepositoriesReponse

	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, fmt.Errorf("decode  %q response: %w", searchURL, err)
	}

	return repos.Items, err
}

func (c *Client) GatherLatestRepositoriesStats() ([]byte, error) {
	return nil, nil
}

func (c *Client) getRepoLanguages(repoURL string) ([]byte, error) {
	url := fmt.Sprintf("%s/languages", repoURL)

	return httpRawGet(url)
}

func httpRawGet(rawURL string) ([]byte, error) {
	if _, err := url.Parse(rawURL); err != nil {
		return nil, fmt.Errorf("invalid %q: %w", rawURL, err)
	}

	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("call %q: %w", rawURL, err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("body %q: %w", rawURL, err)
	}

	return body, nil
}

func intoSearchQuery(filters url.Values) string {
	var sb strings.Builder
	sb.WriteString("is:public")

	for name, values := range filters {
		sb.WriteString(fmt.Sprintf("+%s:%s", name, strings.Join(values, ",")))
	}
	return sb.String()
}
