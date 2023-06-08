package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `db:"id" json:"id,omitempty"`
	AccountID    uuid.UUID  `db:"account_id" json:"account_id,omitempty"`
	FirstName    string     `db:"first_name" json:"first_name,omitempty"`
	LastName     string     `db:"last_name" json:"last_name,omitempty"`
	Email        string     `db:"email" json:"email,omitempty"`
	Pwd          string     `db:"pwd" json:"-"`
	Retries      int        `db:"retries" json:"-"`
	IsSuperAdmin bool       `db:"is_super_admin" json:"is_super_admin,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type UserStoreReader interface {
	GetUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUsers(ctx context.Context, accountID uuid.UUID) ([]User, error)
}

type UserStoreWriter interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	UpdatePwd(ctx context.Context, id uuid.UUID, pwd string, accountID uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) error
	ResetRetries(ctx context.Context, id uuid.UUID) error
	IncrementRetries(ctx context.Context, id uuid.UUID) error
}

type UserStore interface {
	UserStoreReader
	UserStoreWriter
}
