// internal/models/repo.go
package models

import "time"

type Repo struct {
	ID                 int    `json:"id" db:"id"`
	Name               string `json:"name" db:"name"`
	FullName           string `json:"full_name" db:"full_name"`
	Owner              string `json:"owner" db:"owner"`
	URL                string `json:"url" db:"url"`
	Visibility         string `json:"visibility" db:"visibility"`
	Private            bool   `json:"private" db:"private"`
	IsFork             bool   `json:"is_fork" db:"is_fork"`
	ForkedFrom         string `json:"forked_from" db:"forked_from"`
	ForkedFromOwner    string `json:"forked_from_owner" db:"forked_from_owner"`
	ForksCount         int    `json:"forks_count" db:"forks_count"`
	ForkedToCount      int    `json:"forked_to_count" db:"forked_to_count"`
	StargazersCount    int    `json:"stargazers_count" db:"stargazers_count"`
	CollaboratorsCount int    `json:"collaborators_count" db:"collaborators_count"`
	CollaboratorsList  string `json:"collaborators_list" db:"collaborators_list"`
	WhoHasAccess       string `json:"who_has_access" db:"who_has_access"`
	Language           string `json:"language" db:"language"`
	SizeKB             int    `json:"size_kb" db:"size_kb"`
	CreatedAt          string `json:"created_at" db:"created_at"`
	UpdatedAt          string `json:"updated_at" db:"updated_at"`
	PushedAt           string `json:"pushed_at" db:"pushed_at"`
	DefaultBranch      string `json:"default_branch" db:"default_branch"`
	Archived           bool   `json:"archived" db:"archived"`
	Disabled           bool   `json:"disabled" db:"disabled"`
	License            string `json:"license" db:"license"`
	Description        string `json:"description" db:"description"`
}

type Collaborator struct {
	ID        int       `json:"id" db:"id"`
	RepoID    int       `json:"repo_id" db:"repo_id"`
	Username  string    `json:"username" db:"username"`
	HasAccess bool      `json:"has_access" db:"has_access"`
	CheckedAt time.Time `json:"checked_at" db:"checked_at"`
}

// Response DTOs
type CollaboratorResponse struct {
	Repo      string    `json:"repo"`
	User      string    `json:"user"`
	HasAccess bool      `json:"hasAccess"`
	CheckedAt time.Time `json:"checkedAt,omitempty"`
}

type SingleRepoResponse struct {
	Repo      string    `json:"repo"`
	Owner     string    `json:"owner"`
	User      string    `json:"user"`
	HasAccess bool      `json:"hasAccess"`
	CheckedAt time.Time `json:"checkedAt"`
}


// // internal/models/repo.go
// package models

// import "time"

// type Repo struct {
// 	ID       int    `json:"id" db:"id"`
// 	Name     string `json:"name" db:"name"`
// 	Owner    string `json:"owner" db:"owner"`
// 	FullName string `json:"full_name" db:"full_name"`
// 	URL      string `json:"url" db:"url"`
// }

// type Collaborator struct {
// 	ID        int       `json:"id" db:"id"`
// 	RepoID    int       `json:"repo_id" db:"repo_id"`
// 	Username  string    `json:"username" db:"username"`
// 	HasAccess bool      `json:"has_access" db:"has_access"`
// 	CheckedAt time.Time `json:"checked_at" db:"checked_at"`
// }

// // Response DTOs
// type CollaboratorResponse struct {
// 	Repo      string    `json:"repo"`
// 	User      string    `json:"user"`
// 	HasAccess bool      `json:"hasAccess"`
// 	CheckedAt time.Time `json:"checkedAt,omitempty"`
// }

// type SingleRepoResponse struct {
// 	Repo      string    `json:"repo"`
// 	Owner     string    `json:"owner"`
// 	User      string    `json:"user"`
// 	HasAccess bool      `json:"hasAccess"`
// 	CheckedAt time.Time `json:"checkedAt"`
// }
