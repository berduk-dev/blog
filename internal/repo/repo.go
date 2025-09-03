package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type Repository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) Repository {
	return Repository{
		db: db,
	}
}

type User struct {
	userID    int
	name      string
	email     string
	isAdmin   bool
	createdAt time.Time
	updatedAt time.Time
}

type CreateUserRequest struct {
	name  string
	email string
}

type Post struct {
	postID    int
	userID    int
	title     string
	body      string
	views     int
	createdAt time.Time
	updatedAt time.Time
}

type CreatePostRequest struct {
	userID int
	title  string
	body   string
}

type Comment struct {
	commentID int
	postID    int
	userID    int
	body      string
	createdAt time.Time
	updatedAt time.Time
}

type CreateCommentRequest struct {
	postID int
	userID int
	body   string
}

func (r *Repository) GetUserByID(ctx context.Context, userID int) (*User, error) {
	var user User

	err := r.db.QueryRow(ctx, "SELECT name, email, is_admin, created_at, updated_at FROM users WHERE user_id = $1", userID).Scan(&user.name, &user.email, &user.isAdmin, &user.createdAt, &user.updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error User r.db.Query: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetPostByID(ctx context.Context, postID int) (*Post, error) {
	var post Post

	err := r.db.QueryRow(ctx, "SELECT user_id, title, body, views, created_at, updated_at FROM posts WHERE user_id = $1", postID).Scan(&post.userID, &post.title, &post.body, &post.views, &post.createdAt, &post.updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error User r.db.Query: %w", err)
	}

	return &post, nil
}

func (r *Repository) GetCommentByID(ctx context.Context, commentID int) (*Comment, error) {
	var comment Comment

	err := r.db.QueryRow(ctx, "SELECT user_id, post_id, body, created_at, updated_at FROM posts WHERE user_id = $1", commentID).Scan(&comment.userID, &comment.postID, &comment.body, &comment.createdAt, &comment.updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error User r.db.Query: %w", err)
	}

	return &comment, nil
}

func (r *Repository) CreateUser(ctx context.Context, params CreateUserRequest) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", params.name, params.email)
	if err != nil {
		return fmt.Errorf("error CreateUser Exec: %w", err)
	}

	return nil
}

func (r *Repository) CreatePost(ctx context.Context, params CreatePostRequest) error {
	_, err := r.db.Exec(ctx, "INSERT INTO posts (user_id, title, body) VALUES ($1, $2, $3)", params.userID, params.title, params.body)
	if err != nil {
		return fmt.Errorf("error CreatePost Exec: %w", err)
	}

	return nil
}

func (r *Repository) CreateComment(ctx context.Context, params CreateCommentRequest) error {
	_, err := r.db.Exec(ctx, "INSERT INTO comments (user_id, title, body) VALUES ($1, $2, $3)", params.userID, params.postID, params.body)
	if err != nil {
		return fmt.Errorf("error CreateComment Exec: %w", err)
	}

	return nil
}
