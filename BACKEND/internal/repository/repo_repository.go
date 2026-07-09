package repository

import (
	"database/sql"

	"github.com/travis2319/RepoWatch/internal/models"
)

type repoRepository struct {
	db *sql.DB
}

func NewRepoRepository(db *sql.DB) RepoRepository {
	return &repoRepository{db: db}
}

const repoColumns = `id, name, owner, full_name, url, visibility, private, is_fork,
	forked_from, forked_from_owner, forks_count, forked_to_count, stargazers_count,
	collaborators_count, collaborators_list, who_has_access, language, size_kb, created_at, updated_at, pushed_at,
	default_branch, archived, disabled, license, description`

func scanRepo(row interface {
	Scan(dest ...interface{}) error
}) (*models.Repo, error) {
	var r models.Repo
	err := row.Scan(
		&r.ID, &r.Name, &r.Owner, &r.FullName, &r.URL, &r.Visibility, &r.Private, &r.IsFork,
		&r.ForkedFrom, &r.ForkedFromOwner, &r.ForksCount, &r.ForkedToCount, &r.StargazersCount,
		&r.CollaboratorsCount, &r.CollaboratorsList, &r.WhoHasAccess, &r.Language, &r.SizeKB, &r.CreatedAt, &r.UpdatedAt, &r.PushedAt,
		&r.DefaultBranch, &r.Archived, &r.Disabled, &r.License, &r.Description,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *repoRepository) GetByName(name, owner string) (*models.Repo, error) {
	row := r.db.QueryRow("SELECT "+repoColumns+" FROM repos WHERE name = ? AND owner = ?", name, owner)
	return scanRepo(row)
}

func (r *repoRepository) Create(repo *models.Repo) error {
	return r.Upsert(repo)
}

// Upsert inserts a new repo row or updates every column on an existing one,
// matched by (name, owner).
func (r *repoRepository) Upsert(repo *models.Repo) error {
	_, err := r.db.Exec(`
		INSERT INTO repos (
			name, owner, full_name, url, visibility, private, is_fork,
			forked_from, forked_from_owner, forks_count, forked_to_count, stargazers_count,
			collaborators_count, collaborators_list, who_has_access, language, size_kb, created_at, updated_at, pushed_at,
			default_branch, archived, disabled, license, description
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(name, owner) DO UPDATE SET
			full_name=excluded.full_name, url=excluded.url, visibility=excluded.visibility,
			private=excluded.private, is_fork=excluded.is_fork, forked_from=excluded.forked_from,
			forked_from_owner=excluded.forked_from_owner, forks_count=excluded.forks_count,
			forked_to_count=excluded.forked_to_count, stargazers_count=excluded.stargazers_count,
			collaborators_count=excluded.collaborators_count, collaborators_list=excluded.collaborators_list,
			who_has_access=excluded.who_has_access, language=excluded.language,
			size_kb=excluded.size_kb, created_at=excluded.created_at, updated_at=excluded.updated_at,
			pushed_at=excluded.pushed_at, default_branch=excluded.default_branch,
			archived=excluded.archived, disabled=excluded.disabled, license=excluded.license,
			description=excluded.description, last_checked=CURRENT_TIMESTAMP
	`,
		repo.Name, repo.Owner, repo.FullName, repo.URL, repo.Visibility, repo.Private, repo.IsFork,
		repo.ForkedFrom, repo.ForkedFromOwner, repo.ForksCount, repo.ForkedToCount, repo.StargazersCount,
		repo.CollaboratorsCount, repo.CollaboratorsList, repo.WhoHasAccess, repo.Language, repo.SizeKB, repo.CreatedAt, repo.UpdatedAt, repo.PushedAt,
		repo.DefaultBranch, repo.Archived, repo.Disabled, repo.License, repo.Description,
	)
	return err
}

func (r *repoRepository) GetAll() ([]*models.Repo, error) {
	rows, err := r.db.Query("SELECT " + repoColumns + " FROM repos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*models.Repo
	for rows.Next() {
		repo, err := scanRepo(rows)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}