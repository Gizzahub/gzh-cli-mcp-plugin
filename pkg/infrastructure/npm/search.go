// Package npm provides infrastructure for searching npm packages.
package npm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// SearchResult represents a search result from npm.
type SearchResult struct {
	Objects []PackageObject `json:"objects"`
	Total   int             `json:"total"`
}

// PackageObject represents a package in search results.
type PackageObject struct {
	Package PackageInfo `json:"package"`
	Score   Score       `json:"score"`
}

// PackageInfo represents npm package information.
type PackageInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Author      *Author  `json:"author,omitempty"`
	Publisher   *Author  `json:"publisher,omitempty"`
	Links       Links    `json:"links"`
}

// Author represents package author.
type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

// Links represents package links.
type Links struct {
	NPM        string `json:"npm"`
	Homepage   string `json:"homepage,omitempty"`
	Repository string `json:"repository,omitempty"`
}

// Score represents package score.
type Score struct {
	Final  float64 `json:"final"`
	Detail Detail  `json:"detail"`
}

// Detail represents score details.
type Detail struct {
	Quality     float64 `json:"quality"`
	Popularity  float64 `json:"popularity"`
	Maintenance float64 `json:"maintenance"`
}

// Client is an npm API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new npm client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    "https://registry.npmjs.org",
	}
}

// Search searches npm for MCP-related packages.
func (c *Client) Search(query string, limit int) (*SearchResult, error) {
	// Search for MCP packages
	searchQuery := fmt.Sprintf("%s mcp", query)
	searchURL := fmt.Sprintf("%s/-/v1/search?text=%s&size=%d",
		c.baseURL, url.QueryEscape(searchQuery), limit)

	resp, err := c.httpClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("search npm: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("npm search failed: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &result, nil
}

// GetPackage gets detailed information about a package.
func (c *Client) GetPackage(name string) (*PackageDetail, error) {
	pkgURL := fmt.Sprintf("%s/%s", c.baseURL, url.PathEscape(name))

	resp, err := c.httpClient.Get(pkgURL)
	if err != nil {
		return nil, fmt.Errorf("get package: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("package '%s' not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("npm get failed: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result PackageDetail
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &result, nil
}

// PackageDetail represents detailed package information.
type PackageDetail struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	DistTags    map[string]string `json:"dist-tags"`
	Versions    map[string]struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"versions"`
	Readme     string `json:"readme"`
	Homepage   string `json:"homepage"`
	Repository struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"repository"`
	Author  *Author `json:"author"`
	License string  `json:"license"`
}

// LatestVersion returns the latest version tag.
func (p *PackageDetail) LatestVersion() string {
	if v, ok := p.DistTags["latest"]; ok {
		return v
	}
	return ""
}
