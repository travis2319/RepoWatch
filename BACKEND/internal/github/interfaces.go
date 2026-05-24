package github

import "net/http"

type Client interface {
	GetRepos(owner string) ([]map[string]interface{}, error)
	CheckCollaborator(owner, repo, user string) (bool, error)
	ValidateToken() (map[string]interface{}, http.Header, error)
}
