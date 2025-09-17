package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/berduk-dev/blog/internal/model"
	"github.com/berduk-dev/blog/internal/model/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) Repository {
	return Repository{
		db: db,
	}
}

type CreateUserRequest struct {
	Name           string
	Email          string
	HashedPassword string
	IsAdmin        bool
}

type UpdateUser struct {
	Name    *string
	Email   *string
	IsAdmin *bool
}

type UpdatePost struct {
	Title int
}

type CreatePostRequest struct {
	UserID int
	Title  string
	Body   string
}

type CreateCommentRequest struct {
	PostID int
	UserID int
	Body   string
}

type UpdateComment struct {
	Body string
}

func (r *Repository) GetUser(ctx context.Context, userID int) (model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, "SELECT id, name, email, is_admin, created_at, updated_at FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("error GetUser QueryRow: %w", err)
	}

	return user, nil
}

func (r *Repository) GetPost(ctx context.Context, postID int) (model.Post, error) {
	var post model.Post

	err := r.db.QueryRow(ctx,
		`SELECT id,
       			user_id,
       			title,
       			body,
       			views,
       			created_at,
       			updated_at FROM posts WHERE id = $1`, postID,
	).Scan(&post.ID, &post.UserID, &post.Title, &post.Body, &post.Views, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return model.Post{}, fmt.Errorf("error GetPost QueryRow: %w", err)
	}

	return post, nil
}

func (r *Repository) GetComment(ctx context.Context, commentID int) (model.Comment, error) {
	var comment model.Comment

	err := r.db.QueryRow(ctx,
		"SELECT id, user_id, post_id, body, created_at, updated_at FROM posts WHERE id = $1",
		commentID).Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return model.Comment{}, fmt.Errorf("error GetComment QueryRow: %w", err)
	}

	return comment, nil
}

func (r *Repository) GetCommentsByPostID(ctx context.Context, postID int) ([]model.Comment, error) {
	rows, err := r.db.Query(ctx, "SELECT id, post_id, user_id, body, created_at, updated_at FROM comments WHERE post_id = $1", postID)
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *Repository) GetPostsByUserID(ctx context.Context, userID int) ([]model.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, title, body, views, created_at, updated_at
		FROM posts WHERE user_id = $1
		ORDER BY created_at desc`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Body, &post.Views, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *Repository) GetPosts(ctx context.Context) ([]model.Post, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, user_id, title, body, views, created_at, updated_at
        FROM posts
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Body, &post.Views, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return posts, nil
}

func (r *Repository) CreateUser(ctx context.Context, user CreateUserRequest) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (name, email, hashed_password, is_admin) VALUES ($1, $2, $3, $4)", user.Name, user.Email, user.HashedPassword, user.IsAdmin)
	if err != nil {
		log.Println("repo.CreateUser: ", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			if strings.Contains(pgErr.ConstraintName, "email") {
				return errs.ErrorEmailAlreadyExists
			}
		}
		return fmt.Errorf("error CreateUser Exec: %w", err)
	}

	return nil
}

func (r *Repository) CreatePost(ctx context.Context, params CreatePostRequest) error {
	_, err := r.db.Exec(ctx, "INSERT INTO posts (user_id, title, body) VALUES ($1, $2, $3)", params.UserID, params.Title, params.Body)
	if err != nil {
		return fmt.Errorf("error CreatePost Exec: %w", err)
	}

	return nil
}

func (r *Repository) CreateComment(ctx context.Context, params CreateCommentRequest) error {
	_, err := r.db.Exec(ctx, "INSERT INTO comments (user_id, post_id, body) VALUES ($1, $2, $3)", params.UserID, params.PostID, params.Body)
	if err != nil {
		return fmt.Errorf("error CreateComment Exec: %w", err)
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("error DeleteUser Exec: %w", err)
	}

	return nil
}

func (r *Repository) DeletePost(ctx context.Context, postID int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM posts WHERE id = $1", postID)
	if err != nil {
		return fmt.Errorf("error DeletePost Exec: %w", err)
	}

	return nil
}

func (r *Repository) DeleteComment(ctx context.Context, commentID int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		return fmt.Errorf("error DeleteComment Exec: %w", err)
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, userID int, user UpdateUser) (model.User, error) {
	updatedUser := model.User{}
	err := r.db.QueryRow(
		ctx,
		`UPDATE users
			set name=COALESCE($1, name),
			    email=COALESCE($2, email),
			    is_admin=COALESCE($3, is_admin),
			    updated_at=now() WHERE id=$4
	returning id,
			  name,
			  email,
			  is_admin,
			  created_at,
			  updated_at`,
		user.Name,
		user.Email,
		user.IsAdmin,
		userID,
	).Scan(updatedUser.ID, updatedUser.Name, updatedUser.Email, updatedUser.IsAdmin, updatedUser.CreatedAt, updatedUser.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("error UpdateUser: %w", err)
	}

	return updatedUser, nil
}

func (r *Repository) UpdatePost() {

}

func (r *Repository) UpdateComment(ctx context.Context, commentID int, comment UpdateComment) (model.Comment, error) {
	var updatedComment model.Comment
	err := r.db.QueryRow(
		ctx,
		`UPDATE comments
			set body=$1 WHERE id=$2
	returning id,
	  		  post_id,
			  user_id
			  body
			  created_at
			  updated_at`,
		comment.Body,
		commentID,
	).Scan(updatedComment.ID, updatedComment.PostID, updatedComment.UserID, updatedComment.Body, updatedComment.CreatedAt, updatedComment.UpdatedAt)
	if err != nil {
		return model.Comment{}, fmt.Errorf("error UpdateComment: %w", err)
	}

	return updatedComment, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	err := r.db.QueryRow(
		ctx,
		`SELECT
				id,
				name,
				hashed_password,
				email,
				is_admin,
				created_at,
				updated_at
		   FROM users
		  WHERE email = $1`,
		email).Scan(
		&user.ID,
		&user.Name,
		&user.HashedPassword,
		&user.Email,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("error GetUserByEmail Scan: %w", err)
	}

	return user, nil
}
