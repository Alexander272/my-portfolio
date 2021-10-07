package service

import (
	"context"
	"errors"
	"time"

	"github.com/Alexander272/my-portfolio/internal/domain"
	"github.com/Alexander272/my-portfolio/internal/repository"
	"github.com/Alexander272/my-portfolio/pkg/auth"
	"github.com/Alexander272/my-portfolio/pkg/hash"
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
