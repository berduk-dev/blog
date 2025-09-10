package repo

import (
	"context"
	"fmt"
	"github.com/berduk-dev/blog/internal/model"
	"github.com/jackc/pgx/v5"
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

type CreatePostRequest struct {
	userID int
	title  string
	body   string
}

type UpdatePost struct {
	Title int
}

type CreateCommentRequest struct {
	postID int
	userID int
	body   string
}

type UpdateComment struct {
	body string
}

func (r *Repository) GetUser(ctx context.Context, userID int) (model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, "SELECT id, name, email, is_admin, created_at, updated_at FROM users WHERE user_id = $1", userID).Scan(&user.UserID, &user.Name, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("error User r.db.Query: %w", err)
	}

	return user, nil
}

func (r *Repository) GetPost(ctx context.Context, postID int) (model.Post, error) {
	var post model.Post

	err := r.db.QueryRow(ctx, "SELECT user_id, title, body, views, created_at, updated_at FROM posts WHERE user_id = $1", postID).Scan(&post.UserID, &post.Title, &post.Body, &post.Views, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return model.Post{}, fmt.Errorf("error User r.db.Query: %w", err)
	}

	return post, nil
}

func (r *Repository) GetComment(ctx context.Context, commentID int) (model.Comment, error) {
	var comment model.Comment

	err := r.db.QueryRow(ctx,
		"SELECT user_id, post_id, body, created_at, updated_at FROM posts WHERE user_id = $1",
		commentID).Scan(&comment.UserID, &comment.PostID, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return model.Comment{}, fmt.Errorf("error User r.db.Query: %w", err)
	}

	return comment, nil
}

func (r *Repository) GetPostComments(ctx context.Context, postID string) ([]model.Comment, error) {

	rows, err := r.db.Query(ctx, "SELECT comment_id, post_id, user_id, body, created_at, updated_at FROM comments WHERE post_id = $1", postID)

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt)
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

func (r *Repository) GetUserPosts(ctx context.Context, userID string) ([]model.Post, error) {

	rows, err := r.db.Query(ctx, "SELECT post_id, title, body, views, created_at, updated_at FROM posts WHERE post_id = $1", userID)

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.PostID, &post.Title, &post.Body, post.Views, &post.CreatedAt, &post.UpdatedAt)
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

func (r *Repository) CreateUser(ctx context.Context, user CreateUserRequest) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (name, email, hashed_password, is_admin) VALUES ($1, $2, $3, $4)", user.Name, user.Email, user.HashedPassword, user.IsAdmin)
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

func (r *Repository) DeleteUser(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx, "DELETE * FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("error DeleteUser Exec: %w", err)
	}

	return nil
}

func (r *Repository) DeletePost(ctx context.Context, postID int) error {
	_, err := r.db.Exec(ctx, "DELETE * FROM posts WHERE id = $1", postID)
	if err != nil {
		return fmt.Errorf("error DeletePost Exec: %w", err)
	}

	return nil
}

func (r *Repository) DeleteComment(ctx context.Context, commentID int) error {
	_, err := r.db.Exec(ctx, "DELETE * FROM comments WHERE id = $1", commentID)
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
	).Scan(updatedUser.UserID, updatedUser.Name, updatedUser.Email, updatedUser.IsAdmin, updatedUser.CreatedAt, updatedUser.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("error UpdateUser: %w", err)
	}

	return updatedUser, nil
}

func (r *Repository) UpdatePost() {

}

func (r *Repository) UpdateComment(ctx context.Context, commentID int) (UpdateComment, error) {
	var updatedComment UpdateComment
	err := r.db.QueryRow(ctx, `UPDATE comments set body=$1 WHERE id=$2`, commentID).Scan(updatedComment.body)
	if err != nil {
		return UpdateComment{}, fmt.Errorf("error UpdateComment Exec: %w", err)
	}

	return updatedComment, nil
}
