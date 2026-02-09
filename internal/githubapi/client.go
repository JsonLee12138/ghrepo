package githubapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	clerrors "githubRAGCli/internal/exitcode"
)

// Client wraps GitHub REST API access.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a Client with the given configuration.
func NewClient(baseURL, token string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// UserResult holds data returned from GET /user.
type UserResult struct {
	Login              string
	RateLimitRemaining int
}

// GetAuthenticatedUser calls GET /user and returns identity + rate-limit info.
func (c *Client) GetAuthenticatedUser() (*UserResult, error) {
	url := fmt.Sprintf("%s/user", c.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, clerrors.NewTransport("failed to build request", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, clerrors.ClassifyTransportErr(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		rateLimited := isRateLimited(resp)
		return nil, clerrors.ClassifyHTTP(resp.StatusCode, rateLimited, string(body))
	}

	var user struct {
		Login string `json:"login"`
	}
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, clerrors.NewTransport("failed to parse response", err)
	}

	remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))

	return &UserResult{
		Login:              user.Login,
		RateLimitRemaining: remaining,
	}, nil
}

// GetContents calls GET /repos/{owner}/{repo}/contents/{path} and returns the raw JSON body.
// The response may be a single file object or an array of directory entries.
func (c *Client) GetContents(owner, repo, path, ref string) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", c.BaseURL, owner, repo, path)
	if ref != "" {
		url += "?ref=" + ref
	}
	return c.doGet(url)
}

// TreeEntry represents a single entry from the Git Trees API.
type TreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"` // "blob" or "tree"
	SHA  string `json:"sha"`
	Size int64  `json:"size,omitempty"`
}

// TreeResult holds the response from the Git Trees API.
type TreeResult struct {
	SHA       string      `json:"sha"`
	Tree      []TreeEntry `json:"tree"`
	Truncated bool        `json:"truncated"`
}

// GetTree calls GET /repos/{owner}/{repo}/git/trees/{sha} and returns the tree.
func (c *Client) GetTree(owner, repo, sha string, recursive bool) (*TreeResult, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s", c.BaseURL, owner, repo, sha)
	if recursive {
		url += "?recursive=1"
	}

	raw, err := c.doGet(url)
	if err != nil {
		return nil, err
	}

	var result TreeResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, clerrors.NewTransport("failed to parse tree response", err)
	}
	return &result, nil
}

// doGet performs an authenticated GET request and returns the response body.
func (c *Client) doGet(url string) (json.RawMessage, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, clerrors.NewTransport("failed to build request", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, clerrors.ClassifyTransportErr(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		rateLimited := isRateLimited(resp)
		return nil, clerrors.ClassifyHTTP(resp.StatusCode, rateLimited, string(body))
	}

	return json.RawMessage(body), nil
}

// isRateLimited checks response headers for rate-limit indicators.
func isRateLimited(resp *http.Response) bool {
	remaining := resp.Header.Get("X-RateLimit-Remaining")
	if remaining == "0" {
		return true
	}
	return false
}
