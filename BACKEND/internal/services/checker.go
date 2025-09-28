// internal/services/checker.go
package services

import (
	"fmt"
	"log"
	"time"

	"github.com/travis2319/GITHUB-ACCESS/internal/github"
	"github.com/travis2319/GITHUB-ACCESS/internal/models"
	"github.com/travis2319/GITHUB-ACCESS/internal/repository"
)

type CheckerService struct {
	repoRepo     repository.RepoRepository
	collabRepo   repository.CollaboratorRepository
	githubClient github.Client
}

type CheckerServiceInterface interface {
	CheckCollaborators(owner, user string) ([]*models.CollaboratorResponse, error)
	CheckSingleRepo(owner, repo, user string) (*models.SingleRepoResponse, error)
}

func NewCheckerService(repoRepo repository.RepoRepository, collabRepo repository.CollaboratorRepository, githubClient github.Client) CheckerServiceInterface {
	return &CheckerService{
		repoRepo:     repoRepo,
		collabRepo:   collabRepo,
		githubClient: githubClient,
	}
}

func (s *CheckerService) CheckCollaborators(owner, user string) ([]*models.CollaboratorResponse, error) {
	repos, err := s.githubClient.GetRepos(owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get repos: %w", err)
	}

	var results []*models.CollaboratorResponse

	for _, repo := range repos {
		name := repo["name"].(string)

		// Check collaborator access
		hasAccess, err := s.githubClient.CheckCollaborator(owner, name, user)
		if err != nil {
			log.Printf("Error checking collaborator %s for repo %s: %v", user, name, err)
			hasAccess = false // default to false on error
		}

		// Get or create repo record
		// repoModel, err := s.repoRepo.GetByName(name, owner)
		// if err != nil {
		// 	// Create new repo if not exists
		// 	repoModel = &models.Repo{
		// 		Name:     name,
		// 		Owner:    owner,
		// 		FullName: fmt.Sprintf("%s/%s", owner, name),
		// 	}
		// 	if createErr := s.repoRepo.Create(repoModel); createErr != nil {
		// 		log.Printf("Error creating repo record: %v", createErr)
		// 	}
		// }

		checkedAt := time.Now()

		// Save/update collaborator record
		// if repoModel.ID != 0 {
		// 	if err := s.collabRepo.CreateOrUpdate(repoModel.ID, user, hasAccess, checkedAt); err != nil {
		// 		log.Printf("Error saving collaborator record: %v", err)
		// 	}
		// }

		results = append(results, &models.CollaboratorResponse{
			Repo:      name,
			User:      user,
			HasAccess: hasAccess,
			CheckedAt: checkedAt,
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

	// Get or create repo record
	// repoModel, err := s.repoRepo.GetByName(repo, owner)
	// if err != nil {
	// 	// Create new repo if not exists
	// 	repoModel = &models.Repo{
	// 		Name:     repo,
	// 		Owner:    owner,
	// 		FullName: fmt.Sprintf("%s/%s", owner, repo),
	// 	}
	// 	if createErr := s.repoRepo.Create(repoModel); createErr != nil {
	// 		log.Printf("Error creating repo record: %v", createErr)
	// 	}
	// }

	checkedAt := time.Now()

	// Save/update collaborator record
	// if repoModel.ID != 0 {
	// 	if err := s.collabRepo.CreateOrUpdate(repoModel.ID, user, hasAccess, checkedAt); err != nil {
	// 		log.Printf("Error saving collaborator record: %v", err)
	// 	}
	// }

	result := &models.SingleRepoResponse{
		Repo:      repo,
		Owner:     owner,
		User:      user,
		HasAccess: hasAccess,
		CheckedAt: checkedAt,
	}

	log.Printf("Result: %+v", result)
	return result, nil
}
