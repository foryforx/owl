package pg

import (
	"context"
	"database/sql"

	"github.com/foryforx/owl/owl-auth/internal/db"
	"github.com/foryforx/owl/owl-auth/internal/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func NewUserStore(db db.DBTx) *UserStore {
	return &UserStore{db: db}
}

type UserStore struct {
	db db.DBTx
}

func (u *UserStore) GetUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (domain.User, error) {
	query := `
		SELECT id, account_id, first_name, last_name, email, pwd, created_at, updated_at, deleted_at, retries, is_super_admin
		FROM users
		WHERE id = $1 AND account_id = $2;
	`

	var user domain.User
	err := GetContext(ctx, u.db, &user, query, id, accountID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserStore) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, account_id, first_name, last_name, email, pwd, created_at, updated_at, deleted_at, retries, is_super_admin
		FROM users
		WHERE email = $1;
	`

	var user domain.User
	err := GetContext(ctx, u.db, &user, query, email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserStore) GetUsers(ctx context.Context, accountID uuid.UUID) ([]domain.User, error) {
	query := `
		SELECT id, account_id, first_name, last_name, email, pwd, created_at, updated_at, deleted_at, retries, is_super_admin
		FROM users
		WHERE account_id = $1;
	`

	var users []domain.User
	if err := SelectContext(ctx, u.db, &users, query, accountID); err != nil {
		return users, errors.Wrapf(err, "UserStore.GetUsers(%s)", accountID)
	}
	return users, nil
}

func (u *UserStore) CreateUser(ctx context.Context, user *domain.User) error {
	tx, err := u.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	log.Infoln("Creating user", user)
	query := `INSERT INTO users(id, account_id, first_name, last_name, email, pwd, is_super_admin) 
						VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.ExecContext(ctx, query,
		user.ID,
		user.AccountID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Pwd,
		user.IsSuperAdmin,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *UserStore) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `UPDATE users
						SET 
							first_name = $1,
							last_name = $2,
							email = $3,
							is_super_admin = $6
						WHERE
							id = $4
							AND account_id = $5`

	_, err := u.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.ID,
		user.AccountID,
		user.IsSuperAdmin,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) UpdatePwd(ctx context.Context, id uuid.UUID, pwd string, accountID uuid.UUID) error {
	query := `UPDATE users
						SET 
							pwd = $1
						WHERE
							id = $2
							AND account_id = $3`

	_, err := u.db.ExecContext(ctx, query, pwd, id, accountID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) DeleteUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) error {
	query := `UPDATE users
						SET deleted_at = now()
						WHERE
							id = $1
							AND account_id = $2`

	_, err := u.db.ExecContext(ctx, query, id, accountID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) ResetRetries(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users
						SET 
							retries = 0
						WHERE
							id = $1`

	_, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) IncrementRetries(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users
						SET 
							retries = retries + 1
						WHERE
							id = $1`

	_, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
