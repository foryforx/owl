package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID  `db:"id" json:"id,omitempty"`
	Name      string     `db:"name" json:"name,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type AccountStoreReader interface {
	GetAccount(ctx context.Context, id uuid.UUID) (Account, error)
	GetAccounts(ctx context.Context) ([]Account, error)
	GetAccountByName(ctx context.Context, name string) (Account, error)
}

type AccountStoreWriter interface {
	CreateAccount(ctx context.Context, user *Account) error
	UpdateAccount(ctx context.Context, user *Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

type AccountStore interface {
	AccountStoreReader
	AccountStoreWriter
}
