// internal/services/checker.go
package services

import (
	"bytes"
	"fmt"
	"log"
	"time"
	"github.com/travis2319/RepoWatch/internal/github"
	"github.com/travis2319/RepoWatch/internal/models"
	"github.com/travis2319/RepoWatch/internal/repository"
	"github.com/xuri/excelize/v2"
)

type CheckerService struct {
	repoRepo     repository.RepoRepository
	collabRepo   repository.CollaboratorRepository
	githubClient github.Client
}

type GitHubValidationResponse struct {
	Valid      bool                   `json:"valid"`
	User       map[string]interface{} `json:"user,omitempty"`
	RateLimit  map[string]interface{} `json:"rate_limit,omitempty"`
	TokenExpiry string                `json:"token_expiry,omitempty"`
}

type CheckerServiceInterface interface {
	CheckCollaborators(owner, user string) ([]*models.CollaboratorResponse, error)
	CheckSingleRepo(owner, repo, user string) (*models.SingleRepoResponse, error)
	GetOwnerRepos(owner string) ([]map[string]interface{}, error)
	ValidateGitHubToken() (*GitHubValidationResponse, error)
	LoadRepos(owner string) ([]*models.Repo, error)
	LoadCollaborators(owner, repo string) ([]*models.Collaborator, error)
	ExportToExcel() (*bytes.Buffer, error)
}

func NewCheckerService(repoRepo repository.RepoRepository, collabRepo repository.CollaboratorRepository, githubClient github.Client) CheckerServiceInterface {
	return &CheckerService{
		repoRepo:     repoRepo,
		collabRepo:   collabRepo,
		githubClient: githubClient,
	}
}

func (s *CheckerService) GetOwnerRepos(owner string) ([]map[string]interface{}, error) {
	repos, err := s.githubClient.GetRepos(owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get repos for owner %s: %w", owner, err)
	}
	return repos, nil
}

func (s *CheckerService) CheckCollaborators(owner, user string) ([]*models.CollaboratorResponse, error) {
	repos, err := s.githubClient.GetRepos(owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get repos: %w", err)
	}

	var results []*models.CollaboratorResponse

	for _, repo := range repos {
		name, ok := repo["name"].(string)
		if !ok || name == "" {
			continue
		}

		hasAccess, err := s.githubClient.CheckCollaborator(owner, name, user)
		if err != nil {
			log.Printf("Error checking collaborator %s for repo %s: %v", user, name, err)
			hasAccess = false
		}

		results = append(results, &models.CollaboratorResponse{
			Repo:      name,
			User:      user,
			HasAccess: hasAccess,
			CheckedAt: time.Now(),
		})
	}

	return results, nil
}

func (s *CheckerService) CheckSingleRepo(owner, repo, user string) (*models.SingleRepoResponse, error) {
	log.Printf("Checking access for user=%s on repo=%s/%s", user, owner, repo)

	hasAccess, err := s.githubClient.CheckCollaborator(owner, repo, user)
	if err != nil {
		return nil, fmt.Errorf("failed to check collaborator: %w", err)
	}

	result := &models.SingleRepoResponse{
		Repo:      repo,
		Owner:     owner,
		User:      user,
		HasAccess: hasAccess,
		CheckedAt: time.Now(),
	}

	log.Printf("Result: %+v", result)
	return result, nil
}

func (s *CheckerService) ValidateGitHubToken() (*GitHubValidationResponse, error) {
	user, headers, err := s.githubClient.ValidateToken()
	if err != nil {
		return &GitHubValidationResponse{
			Valid: false,
		}, err
	}

	response := &GitHubValidationResponse{
		Valid: true,
		User: map[string]interface{}{
			"login":      user["login"],
			"id":         user["id"],
			"name":       user["name"],
			"avatar_url": user["avatar_url"],
			"html_url":   user["html_url"],
		},
		RateLimit: map[string]interface{}{
			"limit":     headers.Get("X-RateLimit-Limit"),
			"remaining": headers.Get("X-RateLimit-Remaining"),
			"used":      headers.Get("X-RateLimit-Used"),
			"reset":     headers.Get("X-RateLimit-Reset"),
			"resource":  headers.Get("X-RateLimit-Resource"),
		},
		TokenExpiry: headers.Get("github-authentication-token-expiration"),
	}

	return response, nil
}

// LoadRepos fetches repos from GitHub for an owner and persists any new
// ones to the local DB, returning the full saved list.
func (s *CheckerService) LoadRepos(owner string) ([]*models.Repo, error) {
	rawRepos, err := s.githubClient.GetRepos(owner)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos from github: %w", err)
	}

	var saved []*models.Repo
	for _, r := range rawRepos {
		name, _ := r["name"].(string)
		fullName, _ := r["full_name"].(string)
		htmlURL, _ := r["html_url"].(string)
		if name == "" {
			continue
		}
		if htmlURL == "" {
			htmlURL = fmt.Sprintf("https://github.com/%s/%s", owner, name)
		}

		if existing, err := s.repoRepo.GetByName(name, owner); err == nil && existing != nil {
			saved = append(saved, existing)
			continue
		}

		repo := &models.Repo{Name: name, Owner: owner, FullName: fullName, URL: htmlURL}
		if err := s.repoRepo.Create(repo); err != nil {
			log.Printf("failed to save repo %s: %v", name, err)
			continue
		}

		stored, err := s.repoRepo.GetByName(name, owner)
		if err != nil {
			log.Printf("failed to reload saved repo %s: %v", name, err)
			continue
		}
		saved = append(saved, stored)
	}

	return saved, nil
}

// LoadCollaborators fetches collaborators for a single repo from GitHub and
// persists them. Creates the repo record first if it doesn't exist yet
// (so you can call this without calling LoadRepos first).
func (s *CheckerService) LoadCollaborators(owner, repoName string) ([]*models.Collaborator, error) {
	repo, err := s.repoRepo.GetByName(repoName, owner)
	if err != nil {
		// repo = &models.Repo{Name: repoName, Owner: owner, FullName: owner + "/" + repoName}
		repo = &models.Repo{
			Name:     repoName,
			Owner:    owner,
			FullName: owner + "/" + repoName,
			URL:      fmt.Sprintf("https://github.com/%s/%s", owner, repoName),
		}		
		if err := s.repoRepo.Create(repo); err != nil {
			return nil, fmt.Errorf("failed to create repo record: %w", err)
		}
		repo, err = s.repoRepo.GetByName(repoName, owner)
		if err != nil {
			return nil, fmt.Errorf("failed to load repo after creation: %w", err)
		}
	}

	rawCollabs, err := s.githubClient.GetCollaborators(owner, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collaborators: %w", err)
	}

	now := time.Now()
	for _, c := range rawCollabs {
		username, ok := c["login"].(string)
		if !ok || username == "" {
			continue
		}
		if err := s.collabRepo.CreateOrUpdate(repo.ID, username, true, now); err != nil {
			log.Printf("failed to save collaborator %s: %v", username, err)
		}
	}

	return s.collabRepo.GetByRepo(repo.ID)
}

// ExportToExcel builds an .xlsx of every saved repo and its collaborators.
func (s *CheckerService) ExportToExcel() (*bytes.Buffer, error) {
	repos, err := s.repoRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load repos: %w", err)
	}

	f := excelize.NewFile()
	sheet := "Collaborators"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Repo", "Owner", "Full Name", "Username", "Has Access", "Checked At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	row := 2
	for _, repo := range repos {
		collabs, err := s.collabRepo.GetByRepo(repo.ID)
		if err != nil {
			log.Printf("failed to load collaborators for repo %s: %v", repo.Name, err)
			continue
		}

		if len(collabs) == 0 {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), repo.Name)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), repo.Owner)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), repo.FullName)
			row++
			continue
		}

		for _, c := range collabs {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), repo.Name)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), repo.Owner)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), repo.FullName)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), c.Username)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), c.HasAccess)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), c.CheckedAt.Format(time.RFC3339))
			row++
		}
	}

	f.SetColWidth(sheet, "A", "C", 25)
	f.SetColWidth(sheet, "D", "F", 20)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write excel buffer: %w", err)
	}
	return buf, nil
}