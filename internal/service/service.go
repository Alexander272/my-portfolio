package service

import "github.com/Alexander272/my-portfolio/internal/repository"

type Services struct{}

type Deps struct {
	Repos *repository.Repositories
}

func NewServices(deps Deps) *Services {
	return &Services{}
}
