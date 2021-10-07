package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Alexander272/my-portfolio/internal/repository"
	"github.com/Alexander272/my-portfolio/pkg/auth"
	"github.com/Alexander272/my-portfolio/pkg/hash"
)

type AuthService struct {
	repoUsers       repository.Users
	repoAuth        repository.Auth
	tokenManager    auth.TokenManager
	hasher          hash.PasswordHasher
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	domain          string
}

func NewAuthService(repoUsers repository.Users, repoAuth repository.Auth, tokenManager auth.TokenManager, hasher hash.PasswordHasher,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration, domain string) *AuthService {
	return &AuthService{
		repoUsers:       repoUsers,
		repoAuth:        repoAuth,
		tokenManager:    tokenManager,
		hasher:          hasher,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		domain:          domain,
	}
}

func (s *AuthService) SignIn(ctx context.Context, input SignInInput, ua, ip string) (*http.Cookie, *Token, error) {
	user, err := s.repoUsers.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}
	if ok := s.hasher.CheckPasswordHash(input.Password, user.Password); !ok {
		return nil, nil, errors.New("invalid credentials")
	}

	accessToken, err := s.tokenManager.NewJWT(user.Id, user.Email, user.Role, s.accessTokenTTL)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, nil, err
	}

	if err := s.repoAuth.CreateSession(refreshToken, repository.RedisData{
		UserId: user.Id,
		Email:  user.Email,
		Role:   user.Role,
		Ua:     ua,
		Ip:     ip,
		Exp:    s.refreshTokenTTL,
	}); err != nil {
		return nil, nil, err
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    refreshToken,
		MaxAge:   int(s.refreshTokenTTL.Seconds()),
		Path:     "/",
		Domain:   s.domain,
		Secure:   false,
		HttpOnly: true,
	}

	return cookie, &Token{
		AccessToken: accessToken,
	}, nil
}

func (s *AuthService) SingOut(token string) (*http.Cookie, error) {
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		MaxAge:   1,
		Path:     "/",
		Domain:   s.domain,
		Secure:   false,
		HttpOnly: true,
	}

	err := s.repoAuth.RemoveSession(token)
	if err != nil {
		return cookie, err
	}

	return cookie, nil
}

func (s *AuthService) Refresh(token, ua, ip string) (*Token, *http.Cookie, error) {
	data, err := s.repoAuth.GetDelSession(token)
	if err != nil {
		return nil, nil, err
	}
	if ua != data.Ua || ip != data.Ip {
		return nil, nil, errors.New("invalid data")
	}

	accessToken, err := s.tokenManager.NewJWT(data.UserId, data.Email, data.Role, s.accessTokenTTL)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, nil, err
	}

	if err := s.repoAuth.CreateSession(refreshToken, repository.RedisData{
		UserId: data.UserId,
		Email:  data.Email,
		Role:   data.Role,
		Ua:     ua,
		Ip:     ip,
		Exp:    s.refreshTokenTTL,
	}); err != nil {
		return nil, nil, err
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    refreshToken,
		MaxAge:   int(s.refreshTokenTTL.Seconds()),
		Path:     "/",
		Domain:   s.domain,
		Secure:   false,
		HttpOnly: true,
	}

	return &Token{
		AccessToken: accessToken,
	}, cookie, nil
}

func (s *AuthService) TokenParse(token string) (userId string, role string, err error) {
	claims, err := s.tokenManager.Parse(token)
	if err != nil {
		return "", "", err
	}
	return claims["userId"].(string), claims["role"].(string), err
}
