// internal/services/checker.go
package services

import (
	"fmt"
	"log"
	"time"

	"github.com/travis2319/RepoWatch/internal/github"
	"github.com/travis2319/RepoWatch/internal/models"
	"github.com/travis2319/RepoWatch/internal/repository"
)

type CheckerService struct {
	repoRepo     repository.RepoRepository
	collabRepo   repository.CollaboratorRepository
	githubClient github.Client
}

type CheckerServiceInterface interface {
	CheckCollaborators(owner, user string) ([]*models.CollaboratorResponse, error)
	CheckSingleRepo(owner, repo, user string) (*models.SingleRepoResponse, error)
	GetOwnerRepos(owner string) ([]map[string]interface{}, error)
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