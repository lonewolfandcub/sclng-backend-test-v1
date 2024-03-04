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
	maximumResults            = 100
)

type (
	Client struct {
		searchReposBaseURL string
		maxResults         uint
	}

	repoLanguages struct {
		ID        uint
		Languages []Language
	}
)

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

	var repos GithubRepositoriesReponse

	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, fmt.Errorf("decode  %q response: %w", searchURL, err)
	}

	return repos.Items, err
}

func (c *Client) GatherLatestRepositoriesStats(filters url.Values) ([]StatsRepository, error) {
	repos, err := c.ListLatestRepositories(filters)
	if err != nil {
		return nil, err
	}

	ch := make(chan repoLanguages, len(repos))

	statsReposDict := make(map[uint]StatsRepository)
	for _, repo := range repos {
		statsRepo := StatsRepository{
			ID:       repo.ID,
			FullName: repo.FullName,
			License:  repo.License,
			URL:      repo.URL,
		}

		statsReposDict[repo.ID] = statsRepo

		// Concurrently gather languages informations.
		go func(repoID uint, repoURL string) {
			repoLang := repoLanguages{ID: repoID}
			if languages, err := c.getRepoLanguages(repoURL); err == nil {
				repoLang.Languages = languages
			}
			ch <- repoLang
		}(repo.ID, repo.URL)
	}

	// Complete stats repositories with languages informations.
	statsRepos := make([]StatsRepository, 0, len(statsReposDict))

	for i := range repos {
		_ = i

		repoLang := <-ch
		v := statsReposDict[repoLang.ID]
		v.Languages = repoLang.Languages

		statsRepos = append(statsRepos, v)
	}
	return statsRepos, nil
}

func (c *Client) getRepoLanguages(repoURL string) ([]Language, error) {
	languagesURL := fmt.Sprintf("%s/languages", repoURL)

	body, err := httpRawGet(languagesURL)
	if err != nil {
		return nil, err
	}

	var languagesDict map[string]uint
	if err = json.Unmarshal(body, &languagesDict); err != nil {
		return nil, fmt.Errorf("decode %q response: %w", repoURL, err)
	}

	languages := make([]Language, 0, len(languagesDict))
	for k, v := range languagesDict {
		languages = append(languages, Language{Name: k, Bytes: v})
	}

	return languages, nil
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
