// internal/repository/collaborator_repository.go
package repository

import (
	"database/sql"
	"time"

	"github.com/travis2319/GITHUB-ACCESS/internal/models"
)

type collaboratorRepository struct {
	db *sql.DB
}

func NewCollaboratorRepository(db *sql.DB) CollaboratorRepository {
	return &collaboratorRepository{db: db}
}

func (c *collaboratorRepository) CreateOrUpdate(repoID int, username string, hasAccess bool, checkedAt time.Time) error {
	_, err := c.db.Exec(`
		INSERT OR REPLACE INTO collaborators (repo_id, username, has_access, checked_at)
		VALUES (?, ?, ?, ?)
	`, repoID, username, hasAccess, checkedAt)
	return err
}

func (c *collaboratorRepository) GetByRepo(repoID int) ([]*models.Collaborator, error) {
	rows, err := c.db.Query("SELECT id, repo_id, username, has_access, checked_at FROM collaborators WHERE repo_id = ?", repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collaborators []*models.Collaborator
	for rows.Next() {
		var collab models.Collaborator
		err := rows.Scan(&collab.ID, &collab.RepoID, &collab.Username, &collab.HasAccess, &collab.CheckedAt)
		if err != nil {
			return nil, err
		}
		collaborators = append(collaborators, &collab)
	}
	return collaborators, nil
}

func (c *collaboratorRepository) GetByUser(username string) ([]*models.Collaborator, error) {
	rows, err := c.db.Query("SELECT id, repo_id, username, has_access, checked_at FROM collaborators WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collaborators []*models.Collaborator
	for rows.Next() {
		var collab models.Collaborator
		err := rows.Scan(&collab.ID, &collab.RepoID, &collab.Username, &collab.HasAccess, &collab.CheckedAt)
		if err != nil {
			return nil, err
		}
		collaborators = append(collaborators, &collab)
	}
	return collaborators, nil
}
