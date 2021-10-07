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
	Verify(ctx context.Context, userID primitive.ObjectID, code string) error
	SetSession(ctx context.Context, userID primitive.ObjectID) error
}

type Auth interface {
	CreateSession(token string, data RedisData) error
	GetDelSession(token string) (*RedisData, error)
	RemoveSession(token string) error
}

type Repositories struct {
	Users
	Auth
}

func NewRepositories(db *mongo.Database, client *redis.Client) *Repositories {
	return &Repositories{
		Auth:  NewAuthRepo(client),
		Users: NewUsersRepo(db),
	}
}
