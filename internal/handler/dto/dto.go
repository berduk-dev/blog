package dto

import "time"

type CreatePostReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type CreateCommentReq struct {
	Body string `json:"body"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
