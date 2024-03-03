package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (c *Client) ListLatestRepositories() ([]Repository, error) {
	searchURL := fmt.Sprintf("%s?q=is:public&sort=updated&order=desc&per_page=%d", c.searchReposBaseURL, c.maxResults)

	var repos RepositoriesReponse

	body, err := httpRawGet(searchURL)
	if err != nil {
		return nil, err
	}

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
