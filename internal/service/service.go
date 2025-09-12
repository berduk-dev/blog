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

func New(repo repo.Repository) Service {
	return Service{
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

func (s *Service) CreatePost(ctx context.Context, post model.CreatePostReq) error {
	err := s.repo.CreatePost(ctx, repo.CreatePostRequest{
		Title: post.Title,
		Body:  post.Body,
	})
	if err != nil {
		return fmt.Errorf("error repo.CreatePost: %w", err)
	}

	return nil
}

func (s *Service) CreateComment(ctx context.Context, PostID int, comment model.CreateCommentReq) error {
	err := s.repo.CreateComment(ctx, repo.CreateCommentRequest{
		PostID: PostID,
		Body:   comment.Body,
	})
	if err != nil {
		return fmt.Errorf("error repo.CreateComment: %w", err)
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

func (s *Service) GetPostsByUserID(ctx context.Context, userID int) ([]model.Post, error) {
	posts, err := s.repo.GetPostsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error repo.GetPostsByUserID: %w", err)
	}

	return posts, nil
}

func (s *Service) GetPosts(ctx context.Context) ([]model.Post, error) {
	posts, err := s.repo.GetPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("error repo.GetPosts: %w", err)
	}

	return posts, nil
}

func (s *Service) GetPost(ctx context.Context, postID int) (model.Post, error) {
	posts, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		return model.Post{}, fmt.Errorf("error repo.GetPost: %w", err)
	}

	return posts, nil
}

func (s *Service) GetCommentsByPostID(ctx context.Context, postID int) ([]model.Comment, error) {
	comments, err := s.repo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("error repo.GetCommentsByPostID: %w", err)
	}

	return comments, nil
}
