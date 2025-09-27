// internal/models/repo.go
package models

import "time"

type Repo struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Owner    string `json:"owner" db:"owner"`
	FullName string `json:"full_name" db:"full_name"`
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
