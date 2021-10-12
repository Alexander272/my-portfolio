package service

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/Alexander272/my-portfolio/internal/domain"
	"github.com/Alexander272/my-portfolio/internal/repository"
	"github.com/Alexander272/my-portfolio/pkg/auth"
	"github.com/Alexander272/my-portfolio/pkg/hash"
	"github.com/Alexander272/my-portfolio/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CookieName = "session"

type SignUpInput struct {
	Name     string
	Email    string
	Password string
}
type SignInInput struct {
	Email    string
	Password string
}
type Token struct {
	AccessToken  string
	RefreshToken string
}

type Auth interface {
	SignIn(ctx context.Context, input SignInInput, ua, ip string) (*http.Cookie, *Token, error)
	SingOut(token string) (*http.Cookie, error)
	Refresh(token, ua, ip string) (*Token, *http.Cookie, error)
	TokenParse(token string) (userId string, role string, err error)
}

type User interface {
	SignUp(ctx context.Context, input SignUpInput) error
	GetById(ctx context.Context, userId primitive.ObjectID) (domain.User, error)
	UpdateById(ctx context.Context, userId primitive.ObjectID, user domain.UserUpdate) error
	RemoveById(ctx context.Context, userId primitive.ObjectID) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
}

type File interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, path, filename string) (*domain.File, error)
	Remove(ctx context.Context, path, filename string) error
}

type Services struct {
	Auth
	User
	File
}

type Deps struct {
	Repos                  *repository.Repositories
	StorageProvider        storage.Provider
	Hasher                 hash.PasswordHasher
	TokenManager           auth.TokenManager
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	Domain                 string
	VerificationCodeLength int
}

func NewServices(deps Deps) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repos.Users, deps.Repos.Auth, deps.TokenManager, deps.Hasher, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.Domain),
		User: NewUserService(deps.Repos.Users, deps.TokenManager, deps.Hasher),
		File: NewFileService(deps.StorageProvider),
	}
}
