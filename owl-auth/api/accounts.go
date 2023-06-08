package api

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccountRequest struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type AccountResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type CreateAccountResponse struct {
	ID uuid.UUID `json:"id"`
}

type AccountService interface {
	CreateAccount(ctx context.Context, account *AccountRequest) (uuid.UUID, error)
	GetAccount(ctx context.Context, id uuid.UUID) (*AccountResponse, error)
	GetAccounts(ctx context.Context) ([]*AccountResponse, error)
	UpdateAccount(ctx context.Context, account *AccountRequest) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}
