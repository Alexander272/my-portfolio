package repository

import (
	"context"

	"github.com/Alexander272/my-portfolio/internal/domain"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Verify(ctx context.Context, userId primitive.ObjectID, code string) error
	SetSession(ctx context.Context, userId primitive.ObjectID) error
	GetById(ctx context.Context, userId primitive.ObjectID) (domain.User, error)
	UpdateById(ctx context.Context, userId primitive.ObjectID, user domain.UserUpdate) error
	RemoveById(ctx context.Context, userId primitive.ObjectID) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
}

type Auth interface {
	CreateSession(token string, data RedisData) error
	GetDelSession(token string) (*RedisData, error)
	RemoveSession(token string) error
}

type Projects interface {
	GetProjects(ctx context.Context, userId primitive.ObjectID) ([]domain.ProjectMin, error)
	CreateProject(ctx context.Context, project domain.ProjectInput) error
	GetProjectById(ctx context.Context, projectId primitive.ObjectID) (*domain.Project, error)
	UpdateProject(ctx context.Context, projectId primitive.ObjectID, project domain.SelfProject) error
	RemoveProject(ctx context.Context, projectId primitive.ObjectID) error

	GetDrafts(ctx context.Context, userId primitive.ObjectID) ([]domain.SelfProjectMin, error)
	GetSelfProjectById(ctx context.Context, projectId, userId primitive.ObjectID) (*domain.SelfProject, error)
	GetSelfProjects(ctx context.Context, userId primitive.ObjectID) ([]domain.SelfProjectMin, error)
}

type Repositories struct {
	Users
	Auth
	Projects
}

func NewRepositories(db *mongo.Database, client *redis.Client) *Repositories {
	return &Repositories{
		Auth:     NewAuthRepo(client),
		Users:    NewUsersRepo(db),
		Projects: NewProjectsRepo(db),
	}
}
