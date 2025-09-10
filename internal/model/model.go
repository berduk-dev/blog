package model

import "time"

type User struct {
	UserID    int       `json:"user_id"`
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

type Post struct {
	PostID    int
	UserID    int
	Title     string
	Body      string
	Views     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	CommentID int
	PostID    int
	UserID    int
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
