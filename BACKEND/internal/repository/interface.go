package repository

import (
	"time"

	"github.com/travis2319/RepoWatch/internal/models"
)

type RepoRepository interface {
	GetByName(name, owner string) (*models.Repo, error)
	Create(repo *models.Repo) error
	GetAll() ([]*models.Repo, error)
}

type CollaboratorRepository interface {
	CreateOrUpdate(repoID int, username string, hasAccess bool, checkedAt time.Time) error
	GetByRepo(repoID int) ([]*models.Collaborator, error)
	GetByUser(username string) ([]*models.Collaborator, error)
}
