package api

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	ID           uuid.UUID `json:"id"`
	AccountID    uuid.UUID `json:"accountId"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	Pwd          string    `json:"password"`
	IsSuperAdmin bool      `json:"isSuperAdmin"`
}

type UserResponse struct {
	ID           uuid.UUID  `json:"id"`
	AccountID    uuid.UUID  `json:"accountId"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	Email        string     `json:"email"`
	Password     string     `json:"-"`
	Token        string     `json:"token,omitempty"`
	IsSuperAdmin bool       `json:"isSuperAdmin"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID uuid.UUID `json:"id"`
}

type UserService interface {
	Login(ctx context.Context, email, password string) (*UserResponse, error)
	CreateUser(ctx context.Context, user *UserRequest) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*UserResponse, error)
	GetUsers(ctx context.Context, accountID uuid.UUID) ([]*UserResponse, error)
	UpdateUser(ctx context.Context, user *UserRequest) error
	UpdatePassword(ctx context.Context, id uuid.UUID, pwd string, accountID uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) error
}
