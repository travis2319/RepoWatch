package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type client struct {
	token      string
	httpClient *http.Client
}

func NewClient() Client {
	return &client{
		token:      os.Getenv("GITHUB_TOKEN"),
		httpClient: &http.Client{},
	}
}

func (c *client) fetch(url string, target interface{}) error {
	if c.token == "" {
		log.Println("⚠️  GitHub token not set! API call may fail.")
	}
	log.Printf("Fetching URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("HTTP request error: %v", err)
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("GitHub response status: %s", resp.Status)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("GitHub API returned error status %d: %s", resp.StatusCode, string(body))
	}

	return json.Unmarshal(body, target)
}

func (c *client) GetRepos(owner string) ([]map[string]interface{}, error) {
	log.Printf("Getting repos for owner: %s", owner)

	isTokenOwner := false
	if user, _, err := c.ValidateToken(); err == nil {
		if login, ok := user["login"].(string); ok && login == owner {
			isTokenOwner = true
		}
	}

	var allRepos []map[string]interface{}
	page := 1

	for {
		var url string
		if isTokenOwner {
			url = fmt.Sprintf(
				"https://api.github.com/user/repos?visibility=all&affiliation=owner,collaborator,organization_member&per_page=100&page=%d",
				page,
			)
		} else {
			url = fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", owner, page)
		}

		resp, err := c.fetchRaw(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch repos: %w", err)
		}

		var obj map[string]interface{}
		if err := json.Unmarshal(resp, &obj); err == nil {
			if msg, ok := obj["message"]; ok {
				return nil, fmt.Errorf("GitHub API error: %v", msg)
			}
		}

		var repos []map[string]interface{}
		if err := json.Unmarshal(resp, &repos); err != nil {
			return nil, fmt.Errorf("failed to unmarshal repos: %w", err)
		}
		if len(repos) == 0 {
			break
		}
		allRepos = append(allRepos, repos...)
		if len(repos) < 100 {
			break
		}
		page++
	}

	log.Printf("Fetched %d repos for owner %s (tokenOwner=%v)", len(allRepos), owner, isTokenOwner)
	return allRepos, nil
}

func (c *client) GetRepoDetails(owner, repo string) (map[string]interface{}, error) {
	log.Printf("Getting repo details for %s/%s", owner, repo)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	resp, err := c.fetchRaw(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo details: %w", err)
	}

	var details map[string]interface{}
	if err := json.Unmarshal(resp, &details); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repo details: %w", err)
	}
	if msg, ok := details["message"]; ok {
		return nil, fmt.Errorf("GitHub API error: %v", msg)
	}
	return details, nil
}

// func (c *client) GetRepos(owner string) ([]map[string]interface{}, error) {
// 	log.Printf("Getting repos for owner: %s", owner)
// 	var repos []map[string]interface{}
// 	url := fmt.Sprintf("https://api.github.com/users/%s/repos", owner)

// 	resp, err := c.fetchRaw(url)
// 	if err != nil {
// 		log.Printf("Error fetching repos: %v", err)
// 		return nil, fmt.Errorf("failed to fetch repos: %w", err)
// 	}

// 	// Check if response is an object (error)
// 	var obj map[string]interface{}
// 	if err := json.Unmarshal(resp, &obj); err == nil {
// 		if msg, ok := obj["message"]; ok {
// 			log.Printf("GitHub API returned error: %v", msg)
// 			return nil, fmt.Errorf("GitHub API error: %v", msg)
// 		}
// 	}

// 	// Otherwise unmarshal as array
// 	if err := json.Unmarshal(resp, &repos); err != nil {
// 		log.Printf("Error unmarshalling repos: %v", err)
// 		return nil, fmt.Errorf("failed to unmarshal repos: %w", err)
// 	}

// 	log.Printf("Fetched %d repos for owner %s", len(repos), owner)
// 	return repos, nil
// }

func (c *client) GetCollaborators(owner, repo string) ([]map[string]interface{}, error) {
	log.Printf("Getting collaborators for repo: %s/%s", owner, repo)
	var collaborators []map[string]interface{}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/collaborators", owner, repo)

	resp, err := c.fetchRaw(url)
	if err != nil {
		log.Printf("Error fetching collaborators: %v", err)
		return nil, fmt.Errorf("failed to fetch collaborators: %w", err)
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(resp, &obj); err == nil {
		if msg, ok := obj["message"]; ok {
			log.Printf("GitHub API returned error: %v", msg)
			return nil, fmt.Errorf("GitHub API error: %v", msg)
		}
	}

	if err := json.Unmarshal(resp, &collaborators); err != nil {
		log.Printf("Error unmarshalling collaborators: %v", err)
		return nil, fmt.Errorf("failed to unmarshal collaborators: %w", err)
	}

	log.Printf("Fetched %d collaborators for %s/%s", len(collaborators), owner, repo)
	return collaborators, nil
}

func (c *client) fetchRaw(url string) ([]byte, error) {
	log.Printf("Fetching raw data from URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("HTTP request error: %v", err)
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("GitHub raw response status: %s", resp.Status)
	return io.ReadAll(resp.Body)
}

func (c *client) CheckCollaborator(owner, repo, user string) (bool, error) {
	log.Printf("Checking collaborator: user=%s repo=%s/%s", user, owner, repo)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/collaborators/%s", owner, repo, user)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Error checking collaborator: %v", err)
		return false, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Collaborator check response status: %s", resp.Status)

	if resp.StatusCode == 204 {
		log.Printf("User %s has access to repo %s", user, repo)
		return true, nil
	}

	if resp.StatusCode == 404 {
		log.Printf("User %s does NOT have access to repo %s", user, repo)
		return false, nil
	}

	// For other status codes, consider it an error
	body, _ := io.ReadAll(resp.Body)
	return false, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
}

func (c *client) ValidateToken() (map[string]interface{}, http.Header, error) {
	url := "https://api.github.com/user"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, resp.Header, fmt.Errorf("github api returned status %d: %s", resp.StatusCode, string(body))
	}

	var user map[string]interface{}
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, resp.Header, fmt.Errorf("failed to parse response: %w", err)
	}

	return user, resp.Header, nil
}
