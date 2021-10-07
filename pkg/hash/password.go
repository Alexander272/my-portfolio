package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type BcryptHasher struct {
	minCost     int
	maxCost     int
	defaultCost int
}

func NewBcryptHasher(minCost, defaultCost, maxCost int) *BcryptHasher {
	return &BcryptHasher{
		minCost:     minCost,
		maxCost:     maxCost,
		defaultCost: defaultCost,
	}
}

func (h *BcryptHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.defaultCost)
	return string(bytes), err
}

func (h *BcryptHasher) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
