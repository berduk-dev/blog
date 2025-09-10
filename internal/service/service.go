package service

import (
	"context"
	"fmt"
	"github.com/berduk-dev/blog/internal/model"
	"github.com/berduk-dev/blog/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repo.Repository
}

func New(repo repo.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, user model.CreateUserReq) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error bcrypt.GenerateFromPassword: %w", err)
	}

	err = s.repo.CreateUser(ctx, repo.CreateUserRequest{
		Name:           user.Name,
		HashedPassword: string(hashedPassword),
		Email:          user.Email,
	})
	if err != nil {
		return fmt.Errorf("error repo.CreateUser: %w", err)
	}
	return nil
}

func (s *Service) GetUser(ctx context.Context, userID int) (model.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("error repo.GetUser: %w", err)
	}
	return user, nil
}
