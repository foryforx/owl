package pg

import (
	"context"

	"github.com/foryforx/owl/owl-auth/internal/db"
	"github.com/foryforx/owl/owl-auth/internal/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func NewAccountStore(db db.DBTx) *AccountStore {
	return &AccountStore{db: db}
}

type AccountStore struct {
	db db.DBTx
}

func (u *AccountStore) GetAccount(ctx context.Context, id uuid.UUID) (domain.Account, error) {
	query := `
		SELECT id, name, created_at, updated_at, deleted_at
		FROM accounts
		WHERE id = $1;
	`

	var account domain.Account
	err := GetContext(ctx, u.db, &account, query, id)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (u *AccountStore) GetAccountByName(ctx context.Context, name string) (domain.Account, error) {
	query := `
		SELECT id, name, created_at, updated_at, deleted_at
		FROM accounts
		WHERE name = $1;
	`

	var account domain.Account
	err := GetContext(ctx, u.db, &account, query, name)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (u *AccountStore) GetAccounts(ctx context.Context) ([]domain.Account, error) {
	query := `
		SELECT id, name, created_at, updated_at, deleted_at
		FROM accounts;
	`

	var accounts []domain.Account
	if err := SelectContext(ctx, u.db, &accounts, query); err != nil {
		return accounts, errors.Wrapf(err, "AccountStore.GetAccounts()")
	}
	return accounts, nil
}

func (u *AccountStore) CreateAccount(ctx context.Context, account *domain.Account) error {
	query := `INSERT INTO accounts(id, name) 
						VALUES ($1, $2);`

	_, err := u.db.ExecContext(ctx, query,
		account.ID,
		account.Name,
	)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (u *AccountStore) UpdateAccount(ctx context.Context, account *domain.Account) error {
	query := `UPDATE accounts
						SET 
							name = $1
						WHERE
							id = $2`

	_, err := u.db.ExecContext(ctx, query,
		account.Name,
		account.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *AccountStore) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE accounts
						SET deleted_at = now()
						WHERE
							id = $1`

	_, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
