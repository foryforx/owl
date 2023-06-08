package app

import (
	"context"
	"os"
	"time"

	"github.com/foryforx/owl/owl-auth/api"
	"github.com/foryforx/owl/owl-auth/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	MAX_LOGIN_ATTEMPTS = 3
)

func NewUserService(s domain.UserStore) *UserService {
	return &UserService{userStore: s}
}

type UserService struct {
	userStore domain.UserStore
}

func (u *UserService) Login(ctx context.Context, email, password string) (*api.UserResponse, error) {
	user, err := u.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user.DeletedAt != nil {
		return nil, api.ErrUnauthorized
	}
	if user.Retries >= MAX_LOGIN_ATTEMPTS {
		return nil, api.ErrTooManyAttempts
	}
	isMatch := CheckPasswordHash(password, user.Pwd)
	if !isMatch {
		err = u.userStore.IncrementRetries(ctx, user.ID)
		if err != nil {
			return nil, api.ErrUnauthorized
		}
		return nil, api.ErrUnauthorized
	}

	err = u.userStore.ResetRetries(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "UserService.Login(ID=%s) ResetRetries", user.ID)
	}
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &domain.Claims{
		ID:        user.ID,
		AccountID: user.AccountID,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, errors.Wrapf(err, "UserService.Login(ID=%s) SignedString", user.ID)
	}
	userR := &api.UserResponse{
		ID:           user.ID,
		AccountID:    user.AccountID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Pwd,
		Token:        tokenString,
		IsSuperAdmin: user.IsSuperAdmin,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		DeletedAt:    user.DeletedAt,
	}
	return userR, nil
}

func (u *UserService) CreateUser(ctx context.Context, user *api.UserRequest) (uuid.UUID, error) {
	userM := &domain.User{
		ID:           uuid.New(),
		AccountID:    user.AccountID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Pwd:          user.Pwd,
		IsSuperAdmin: user.IsSuperAdmin,
	}
	var err error
	userM.Pwd, err = HashPassword(user.Pwd)
	if err != nil {
		return uuid.UUID{}, errors.Wrapf(err, "UserService.Create(email=%s)", user.Email)
	}
	existingAccount, err := u.userStore.GetUserByEmail(ctx, user.Email)
	if err != nil && err != domain.ErrNotFound {
		return uuid.UUID{}, errors.Wrapf(err, "UserService.Create(email=%s)", user.Email)
	}
	if existingAccount.Email != "" {
		return uuid.UUID{}, domain.DuplicateEntryError{Err: errors.Errorf("UserService.Create(email=%s)", user.Email)}
	}
	if err := u.userStore.CreateUser(ctx, userM); err != nil {
		return uuid.UUID{}, errors.Wrapf(err, "UserService.Create(email=%s)", user.Email)
	}

	return userM.ID, nil
}

func (u *UserService) GetUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*api.UserResponse, error) {
	user, err := u.userStore.GetUser(ctx, id, accountID)
	if err != nil {
		return nil, errors.Wrapf(err, "UserService.GetUser(id=%s)", id)
	}

	userR := &api.UserResponse{
		ID:           user.ID,
		AccountID:    user.AccountID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Pwd,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		DeletedAt:    user.DeletedAt,
		IsSuperAdmin: user.IsSuperAdmin,
	}
	return userR, nil
}

func (u *UserService) GetUsers(ctx context.Context, accountID uuid.UUID) ([]*api.UserResponse, error) {
	users, err := u.userStore.GetUsers(ctx, accountID)
	if err != nil {
		return nil, errors.Wrapf(err, "UserService.GetUsers(accountID=%s)", accountID)
	}
	usersR := []*api.UserResponse{}
	for _, v := range users {
		usersR = append(usersR, &api.UserResponse{
			ID:           v.ID,
			AccountID:    v.AccountID,
			FirstName:    v.FirstName,
			LastName:     v.LastName,
			Email:        v.Email,
			Password:     v.Pwd,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
			DeletedAt:    v.DeletedAt,
			IsSuperAdmin: v.IsSuperAdmin,
		})
	}
	return usersR, nil
}

func (u *UserService) UpdateUser(ctx context.Context, user *api.UserRequest) error {
	userM := &domain.User{
		ID:           user.ID,
		AccountID:    user.AccountID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		IsSuperAdmin: user.IsSuperAdmin,
	}
	err := u.userStore.UpdateUser(ctx, userM)
	if err != nil {
		return errors.Wrapf(err, "UserService.UpdateUser")
	}

	return nil
}

func (u *UserService) UpdatePassword(ctx context.Context, id uuid.UUID, pwd string, accountID uuid.UUID) error {
	user, err := u.userStore.GetUser(ctx, id, accountID)
	if err != nil {
		return errors.Wrapf(err, "UserService.UpdatePassword(id=%s)-GetUser", id)
	}
	pwdHash, err := HashPassword(user.Pwd)
	if err != nil {
		return errors.Wrapf(err, "UserService.UpdatePassword(id=%s)- Hash", id)
	}
	err = u.userStore.UpdatePwd(ctx, id, pwdHash, accountID)
	if err != nil {
		return errors.Wrapf(err, "UserService.UpdatePwd")
	}
	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) error {
	err := u.userStore.DeleteUser(ctx, id, accountID)
	if err != nil {
		return errors.Wrapf(err, "UserService.DeleteUser")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
