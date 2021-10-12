package service

import (
	"context"
	"errors"

	"github.com/Alexander272/my-portfolio/internal/domain"
	"github.com/Alexander272/my-portfolio/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectService struct {
	repo repository.Projects
}

func NewProjectService(repo repository.Projects) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}

func (s *ProjectService) GetProjects(ctx context.Context, userId primitive.ObjectID) ([]domain.ProjectMin, error) {
	return s.repo.GetProjects(ctx, userId)
}

func (s *ProjectService) GetProjectById(ctx context.Context, projectId primitive.ObjectID) (*domain.Project, error) {
	project, err := s.repo.GetProjectById(ctx, projectId)
	if err != nil {
		return nil, err
	}

	if project.Access == domain.Nobody {
		return nil, errors.New("access forbidden")
	}
	return project, nil
}
