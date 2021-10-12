package service

import (
	"context"
	"errors"
	"time"

	"github.com/Alexander272/my-portfolio/internal/domain"
	"github.com/Alexander272/my-portfolio/internal/repository"
	"github.com/Alexander272/my-portfolio/pkg/auth"
	"github.com/Alexander272/my-portfolio/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo         repository.Users
	tokenManager auth.TokenManager
	hasher       hash.PasswordHasher
}

func NewUserService(repo repository.Users, tokenManager auth.TokenManager, hasher hash.PasswordHasher) *UserService {
	return &UserService{
		repo:         repo,
		tokenManager: tokenManager,
		hasher:       hasher,
	}
}

func (s *UserService) SignUp(ctx context.Context, input SignUpInput) error {
	passwordHash, err := s.hasher.HashPassword(input.Password)
	if err != nil {
		return err
	}
	verificationCode, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return err
	}
	ttl, err := time.ParseDuration("6h")
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         input.Name,
		Password:     passwordHash,
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
		Verification: domain.Verification{
			Code:    verificationCode,
			Expires: time.Now().Add(ttl),
		},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return err
		}
		return err
	}

	return nil
}

func (s *UserService) GetById(ctx context.Context, userId primitive.ObjectID) (domain.User, error) {
	return s.repo.GetById(ctx, userId)
}

func (s *UserService) UpdateById(ctx context.Context, userId primitive.ObjectID, input domain.UserUpdate) error {
	return s.repo.UpdateById(ctx, userId, input)
}

func (s *UserService) RemoveById(ctx context.Context, userId primitive.ObjectID) error {
	return s.repo.RemoveById(ctx, userId)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAllUsers(ctx)
}
