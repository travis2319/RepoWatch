package repository

import (
	"database/sql"

	"github.com/travis2319/GITHUB-ACCESS/internal/models"
)

type repoRepository struct {
	db *sql.DB
}

func NewRepoRepository(db *sql.DB) RepoRepository {
	return &repoRepository{db: db}
}

func (r *repoRepository) GetByName(name, owner string) (*models.Repo, error) {
	var repo models.Repo
	err := r.db.QueryRow("SELECT id, name, owner, full_name FROM repos WHERE name = ? AND owner = ?", name, owner).
		Scan(&repo.ID, &repo.Name, &repo.Owner, &repo.FullName)

	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *repoRepository) Create(repo *models.Repo) error {
	_, err := r.db.Exec("INSERT INTO repos (name, owner, full_name) VALUES (?, ?, ?)",
		repo.Name, repo.Owner, repo.FullName)
	return err
}

func (r *repoRepository) GetAll() ([]*models.Repo, error) {
	rows, err := r.db.Query("SELECT id, name, owner, full_name FROM repos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*models.Repo
	for rows.Next() {
		var repo models.Repo
		err := rows.Scan(&repo.ID, &repo.Name, &repo.Owner, &repo.FullName)
		if err != nil {
			return nil, err
		}
		repos = append(repos, &repo)
	}
	return repos, nil
}
