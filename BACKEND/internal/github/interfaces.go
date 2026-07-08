package github

import "net/http"

type Client interface {
	GetRepos(owner string) ([]map[string]interface{}, error)
	CheckCollaborator(owner, repo, user string) (bool, error)
	GetCollaborators(owner, repo string) ([]map[string]interface{}, error)
	ValidateToken() (map[string]interface{}, http.Header, error)
}
