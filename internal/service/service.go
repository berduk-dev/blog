package service

import (
	"context"
	"fmt"
	handlerDto "github.com/berduk-dev/blog/internal/handler/dto"
	"golang.org/x/crypto/bcrypt"

	"github.com/berduk-dev/blog/internal/cache"
	"github.com/berduk-dev/blog/internal/model"
	serviceDto "github.com/berduk-dev/blog/internal/service/dto"
)

type Repository interface {
	CreatePost(ctx context.Context, params serviceDto.CreatePostRequest) error
	GetPost(ctx context.Context, postID int) (model.Post, error)
	GetPosts(ctx context.Context) ([]model.Post, error)
	GetPostsByUserID(ctx context.Context, userID int) ([]model.Post, error)

	CreateComment(ctx context.Context, params serviceDto.CreateCommentRequest) error
	GetCommentsByPostID(ctx context.Context, postID int) ([]model.Comment, error)

	CreateUser(ctx context.Context, user serviceDto.CreateUserRequest) error
	GetUser(ctx context.Context, userID int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type Service struct {
	repo           Repository
	sessionsManger cache.SessionsCache
}

func New(repo Repository, sessionsManger cache.SessionsCache) *Service {
	return &Service{
		repo:           repo,
		sessionsManger: sessionsManger,
	}
}

func (s *Service) CreateUser(ctx context.Context, user handlerDto.CreateUserReq) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error bcrypt.GenerateFromPassword: %w", err)
	}

	err = s.repo.CreateUser(ctx, serviceDto.CreateUserRequest{
		Name:           user.Name,
		HashedPassword: string(hashedPassword),
		Email:          user.Email,
	})
	if err != nil {
		return fmt.Errorf("error repo.CreateUser: %w", err) // "error repo.CreateUser: email уже существует"
	}

	return nil
}

func (s *Service) CreatePost(ctx context.Context, userID int, post handlerDto.CreatePostReq) error {
	err := s.repo.CreatePost(ctx, serviceDto.CreatePostRequest{
		UserID: userID,
		Title:  post.Title,
		Body:   post.Body,
	})
	if err != nil {
		return fmt.Errorf("error repo.CreatePost: %w", err)
	}

	return nil
}

func (s *Service) CreateComment(ctx context.Context, userID int, PostID int, comment handlerDto.CreateCommentReq) error {
	err := s.repo.CreateComment(ctx, serviceDto.CreateCommentRequest{
		UserID: userID,
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

func (s *Service) AuthenticateUser(ctx context.Context, email string, password string) (model.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, fmt.Errorf("error repo.GetUserByEmail: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return model.User{}, fmt.Errorf("invalid email or password") // TODO: вынести ошибку в пакет
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
