package dto

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

type CreateUserRequest struct {
	Name           string
	Email          string
	HashedPassword string
	IsAdmin        bool
}
