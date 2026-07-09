// internal/services/checker.go
package services

import (
	"bytes"
	"fmt"
	"log"
	"time"
	"strings"
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
		repo := mapRawRepoToModel(owner, r)
		if repo.Name == "" {
			continue
		}

		if repo.IsFork {
			details, err := s.githubClient.GetRepoDetails(owner, repo.Name)
			if err != nil {
				log.Printf("failed to fetch fork parent info for %s: %v", repo.FullName, err)
			} else if parent, ok := details["parent"].(map[string]interface{}); ok {
				repo.ForkedFrom, _ = parent["full_name"].(string)
				if parentOwner, ok := parent["owner"].(map[string]interface{}); ok {
					repo.ForkedFromOwner, _ = parentOwner["login"].(string)
				}
			}
		}

		collabs, err := s.githubClient.GetCollaborators(owner, repo.Name)
		if err != nil {
			log.Printf("failed to fetch collaborators for %s: %v", repo.FullName, err)
			repo.CollaboratorsList = "no collaborator"
			repo.WhoHasAccess = "no collaborator"
		} else {
			repo.CollaboratorsCount = len(collabs)

			var names []string
			var accessDetails []string

			for _, c := range collabs {
				login, ok := c["login"].(string)
				if !ok || login == "" {
					continue
				}
				names = append(names, login)

				level := "read"
				if perms, ok := c["permissions"].(map[string]interface{}); ok {
					if admin, _ := perms["admin"].(bool); admin {
						level = "admin"
					} else if push, _ := perms["push"].(bool); push {
						level = "write"
					}
				}
				accessDetails = append(accessDetails, fmt.Sprintf("%s (%s)", login, level))
			}

			if len(names) == 0 {
				repo.CollaboratorsList = "no collaborator"
				repo.WhoHasAccess = "no collaborator"
			} else {
				repo.CollaboratorsList = strings.Join(names, ", ")
				repo.WhoHasAccess = strings.Join(accessDetails, ", ")
			}
		}

		if err := s.repoRepo.Upsert(repo); err != nil {
			log.Printf("failed to save repo %s: %v", repo.Name, err)
			continue
		}

		stored, err := s.repoRepo.GetByName(repo.Name, owner)
		if err != nil {
			log.Printf("failed to reload saved repo %s: %v", repo.Name, err)
			continue
		}
		saved = append(saved, stored)
	}

	// Second pass: "Forked To Count" = how many repos in THIS tracked batch
	// are forks of a given repo. This is computed locally, not from GitHub,
	// since GitHub doesn't expose a full fork list via the REST API.
	forkCounts := map[string]int{}
	for _, r := range saved {
		if r.ForkedFrom != "" {
			forkCounts[r.ForkedFrom]++
		}
	}
	for _, r := range saved {
		if count, ok := forkCounts[r.FullName]; ok && count != r.ForkedToCount {
			r.ForkedToCount = count
			if err := s.repoRepo.Upsert(r); err != nil {
				log.Printf("failed to update forked_to_count for %s: %v", r.FullName, err)
			}
		}
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
	sheet := "Repositories"
	f.SetSheetName("Sheet1", sheet)

		headers := []string{
			"Repository Name", "Full Name", "Owner", "Visibility", "Private", "Is Fork",
			"Forked From", "Forked From Owner", "Forks Count", "Forked To Count", "Stars",
			"Collaborators Count", "Collaborators", "Who Has Access", "Language", "Size (KB)", "Created At", "Updated At",
			"Pushed At", "Default Branch", "Archived", "Disabled", "License", "Description", "URL",
		}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, repo := range repos {
		row := i + 2
		values := []interface{}{
			repo.Name, repo.FullName, repo.Owner, repo.Visibility, repo.Private, repo.IsFork,
			repo.ForkedFrom, repo.ForkedFromOwner, repo.ForksCount, repo.ForkedToCount, repo.StargazersCount,
			repo.CollaboratorsCount, repo.CollaboratorsList, repo.WhoHasAccess, repo.Language, repo.SizeKB, repo.CreatedAt, repo.UpdatedAt,
			repo.PushedAt, repo.DefaultBranch, repo.Archived, repo.Disabled, repo.License,
			repo.Description, repo.URL,
		}
		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
		}
	}

	f.SetColWidth(sheet, "A", "B", 28)
	f.SetColWidth(sheet, "C", "D", 16)
	f.SetColWidth(sheet, "G", "H", 22)
	f.SetColWidth(sheet, "U", "U", 40)
	f.SetColWidth(sheet, "W", "W", 40)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write excel buffer: %w", err)
	}
	return buf, nil
}

func mapRawRepoToModel(owner string, r map[string]interface{}) *models.Repo {
	repo := &models.Repo{Owner: owner}

	repo.Name, _ = r["name"].(string)
	repo.FullName, _ = r["full_name"].(string)
	repo.URL, _ = r["html_url"].(string)
	repo.Visibility, _ = r["visibility"].(string)
	repo.Private, _ = r["private"].(bool)
	repo.IsFork, _ = r["fork"].(bool)
	repo.Language, _ = r["language"].(string)
	repo.CreatedAt, _ = r["created_at"].(string)
	repo.UpdatedAt, _ = r["updated_at"].(string)
	repo.PushedAt, _ = r["pushed_at"].(string)
	repo.DefaultBranch, _ = r["default_branch"].(string)
	repo.Archived, _ = r["archived"].(bool)
	repo.Disabled, _ = r["disabled"].(bool)
	repo.Description, _ = r["description"].(string)

	if v, ok := r["forks_count"].(float64); ok {
		repo.ForksCount = int(v)
	}
	if v, ok := r["stargazers_count"].(float64); ok {
		repo.StargazersCount = int(v)
	}
	if v, ok := r["size"].(float64); ok {
		repo.SizeKB = int(v)
	}
	if lic, ok := r["license"].(map[string]interface{}); ok && lic != nil {
		repo.License, _ = lic["name"].(string)
	}
	if repo.URL == "" && repo.FullName != "" {
		repo.URL = "https://github.com/" + repo.FullName
	}

	return repo
}