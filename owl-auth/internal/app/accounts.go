package app

import (
	"context"

	"github.com/foryforx/owl/owl-auth/api"
	"github.com/foryforx/owl/owl-auth/internal/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func NewAccountService(s domain.AccountStore) *AccountService {
	return &AccountService{AccountStore: s}
}

type AccountService struct {
	AccountStore domain.AccountStore
}

func (u *AccountService) CreateAccount(ctx context.Context, account *api.AccountRequest) (uuid.UUID, error) {
	accountM := &domain.Account{
		ID:   uuid.New(),
		Name: account.Name,
	}

	existingAccount, err := u.AccountStore.GetAccountByName(ctx, account.Name)
	if err != nil && err != domain.ErrNotFound {
		return uuid.UUID{}, errors.Wrapf(err, "AccountService.Create(name=%s)", account.Name)
	}
	if existingAccount.Name != "" {
		return uuid.UUID{}, domain.DuplicateEntryError{Err: errors.Errorf("AccountService.Create(name=%s)", account.Name)}
	}
	if err := u.AccountStore.CreateAccount(ctx, accountM); err != nil {
		return uuid.UUID{}, errors.Wrapf(err, "AccountService.Create(name=%s)", account.Name)
	}

	return accountM.ID, nil
}

func (u *AccountService) GetAccount(ctx context.Context, id uuid.UUID) (*api.AccountResponse, error) {
	account, err := u.AccountStore.GetAccount(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "AccountService.GetAccount(id=%s)", id)
	}
	acc := &api.AccountResponse{
		ID:        account.ID,
		Name:      account.Name,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
		DeletedAt: account.DeletedAt,
	}

	return acc, nil
}

func (u *AccountService) GetAccounts(ctx context.Context) ([]*api.AccountResponse, error) {
	accounts, err := u.AccountStore.GetAccounts(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "AccountService.GetAccount: %w")
	}
	accs := []*api.AccountResponse{}
	for _, a := range accounts {
		accs = append(accs, &api.AccountResponse{
			ID:        a.ID,
			Name:      a.Name,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			DeletedAt: a.DeletedAt,
		})
	}

	return accs, nil
}

func (u *AccountService) UpdateAccount(ctx context.Context, account *api.AccountRequest) error {
	accountM := &domain.Account{
		ID:   account.ID,
		Name: account.Name,
	}
	err := u.AccountStore.UpdateAccount(ctx, accountM)
	if err != nil {
		return errors.Wrapf(err, "AccountService.UpdateAccount: %w")
	}

	return nil
}

func (u *AccountService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	err := u.AccountStore.DeleteAccount(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "AccountService.DeleteAccount: %w")
	}

	return nil
}
