package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/foryforx/owl/owl-auth/api"
	"github.com/foryforx/owl/owl-auth/internal/domain"
	"github.com/foryforx/owl/owl-auth/internal/transport/httpapi/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserRouter struct {
	UserService    api.UserService
	AccountService api.AccountService
}

func (u *UserRouter) Routes() []Route {
	return []Route{
		{http.MethodPost, "/login", LoginHandler(u.UserService), true},
		{http.MethodPost, "/accounts/{account_id}/users", CreateUserHandler(u.UserService, u.AccountService), true},
		{http.MethodGet, "/accounts/{account_id}/users", GetUsersHandler(u.UserService, u.AccountService), false},
		{http.MethodGet, "/accounts/{account_id}/users/{id}", GetUserHandler(u.UserService, u.AccountService), false},
		{http.MethodPut, "/accounts/{account_id}/users/{id}", UpdateUserHandler(u.UserService, u.AccountService), false},
		{http.MethodPatch, "/accounts/{account_id}/users/{id}", UpdatePasswordHandler(u.UserService, u.AccountService), false},
		{http.MethodDelete, "/accounts/{account_id}/users/{id}", DeleteUserHandler(u.UserService, u.AccountService), false},
	}
}

type usersGetter interface {
	GetUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*api.UserResponse, error)
	GetUsers(ctx context.Context, accountID uuid.UUID) ([]*api.UserResponse, error)
	Login(ctx context.Context, email, password string) (*api.UserResponse, error)
}

type usersCreator interface {
	CreateUser(ctx context.Context, user *api.UserRequest) (uuid.UUID, error)
}

type usersUpdator interface {
	UpdateUser(ctx context.Context, user *api.UserRequest) error
	UpdatePassword(ctx context.Context, id uuid.UUID, pwd string, accountID uuid.UUID) error
}

type usersDeleter interface {
	DeleteUser(ctx context.Context, id uuid.UUID, accountID uuid.UUID) error
}

func CreateUserHandler(userSvc usersCreator, acctSvc accountsCreator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// currentUser, ok := r.Context().Value(domain.CKey("user")).(*domain.User)
		// if !ok {
		// 	response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
		// 	return
		// }
		var user api.UserRequest
		accountIDStr := mux.Vars(r)["account_id"]
		if accountIDStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		accountID, err := uuid.Parse(accountIDStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.RespondError(r.Context(), w, err, "invalid request", http.StatusBadRequest)
			return
		}
		// if currentUser.AccountID != accountID {
		// 	response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
		// 	return
		// }

		user.Pwd = "abcde"
		if strings.TrimSpace(user.Email) != "" && !strings.Contains(strings.TrimSpace(user.Email), "@") ||
			strings.TrimSpace(user.FirstName) == "" || strings.TrimSpace(user.LastName) == "" ||
			user.Pwd == "" {
			response.RespondError(r.Context(), w, errors.New("invalid data"), "invalid data", http.StatusBadRequest)
			return
		}
		// if !currentUser.IsSuperAdmin {
		// 	response.RespondError(r.Context(), w, errors.New("invalid access"), "invalid access", http.StatusUnauthorized)
		// 	return
		// }
		user.AccountID = accountID
		id, err := userSvc.CreateUser(r.Context(), &user)
		switch {
		case errors.As(err, &domain.DuplicateEntryError{}):
			response.RespondError(r.Context(), w, err, "rmail already exists", http.StatusConflict)
			return
		case err != nil:
			response.RespondError(r.Context(), w, err, "registration failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, api.CreateUserResponse{ID: id}, http.StatusCreated)
	}
}

func UpdateUserHandler(userSvc usersUpdator, acctSvc accountsCreator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := r.Context().Value(domain.CKey("user")).(*domain.User)
		if !ok {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		var user api.UserRequest
		idStr := mux.Vars(r)["id"]
		if idStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		accountIDStr := mux.Vars(r)["account_id"]
		if accountIDStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		accountID, err := uuid.Parse(accountIDStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		if currentUser.ID != id {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		user.AccountID = accountID
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.RespondError(r.Context(), w, err, "invalid request", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(user.Email) != "" && !strings.Contains(strings.TrimSpace(user.Email), "@") || strings.TrimSpace(user.FirstName) == "" || strings.TrimSpace(user.LastName) == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "invalid request", http.StatusBadRequest)
			return
		}
		user.ID = id
		if err := userSvc.UpdateUser(r.Context(), &user); err != nil {
			response.RespondError(r.Context(), w, err, "update failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, nil, http.StatusOK)
	}
}

func GetUserHandler(userSvc usersGetter, acctSvc accountsCreator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := r.Context().Value(domain.CKey("user")).(*domain.User)
		if !ok {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		idStr := mux.Vars(r)["id"]
		if idStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		accountIDStr := mux.Vars(r)["account_id"]
		if accountIDStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		accountID, err := uuid.Parse(accountIDStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		if currentUser.ID != id {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		user, err := userSvc.GetUser(r.Context(), id, accountID)
		switch {
		case errors.Is(err, domain.ErrNotFound):
			response.RespondError(r.Context(), w, err, "user not found", http.StatusNotFound)
			return
		case err != nil:
			response.RespondError(r.Context(), w, err, "get users failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, user, http.StatusOK)
	}
}

func GetUsersHandler(userSvc usersGetter, acctSvc accountsCreator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := r.Context().Value(domain.CKey("user")).(*domain.User)
		if !ok {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		accountIDStr := mux.Vars(r)["account_id"]
		if accountIDStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		accountID, err := uuid.Parse(accountIDStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		if currentUser.AccountID != accountID {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		if !currentUser.IsSuperAdmin {
			response.RespondError(r.Context(), w, errors.New("invalid access"), "invalid access", http.StatusUnauthorized)
			return
		}
		users, err := userSvc.GetUsers(r.Context(), accountID)
		if err != nil {
			response.RespondError(r.Context(), w, err, "get users failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, users, http.StatusOK)
	}
}

func DeleteUserHandler(userSvc usersDeleter, acctSvc accountsCreator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := r.Context().Value(domain.CKey("user")).(*domain.User)
		if !ok {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		idStr := mux.Vars(r)["id"]
		if idStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		accountIDStr := mux.Vars(r)["account_id"]
		if accountIDStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		accountID, err := uuid.Parse(accountIDStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		if currentUser.AccountID != accountID {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		if !currentUser.IsSuperAdmin {
			response.RespondError(r.Context(), w, errors.New("invalid access"), "invalid access", http.StatusUnauthorized)
			return
		}
		if err := userSvc.DeleteUser(r.Context(), id, accountID); err != nil {
			response.RespondError(r.Context(), w, err, "delete failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, nil, http.StatusOK)
	}
}

func UpdatePasswordHandler(userSvc usersUpdator, acctSvc accountsCreator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := r.Context().Value(domain.CKey("user")).(*domain.User)
		if !ok {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		var user api.UserRequest
		idStr := mux.Vars(r)["id"]
		if idStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		accountIDStr := mux.Vars(r)["account_id"]
		if accountIDStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "user doesnt exists", http.StatusNotFound)
			return
		}
		accountID, err := uuid.Parse(accountIDStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "user doesnt exists", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.RespondError(r.Context(), w, err, "invalid request", http.StatusBadRequest)
			return
		}
		if user.Pwd == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "invalid request", http.StatusBadRequest)
			return
		}
		if currentUser.ID != id {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}

		user.AccountID = accountID
		user.ID = id
		if err := userSvc.UpdatePassword(r.Context(), id, user.Pwd, user.AccountID); err != nil {
			response.RespondError(r.Context(), w, err, "update failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, nil, http.StatusOK)
	}
}

func LoginHandler(userSvc usersGetter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user api.UserLoginRequest

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.RespondError(r.Context(), w, err, "Invalid request", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(user.Email) == "" || user.Password == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "Invalid request", http.StatusBadRequest)
			return
		}

		userRes, err := userSvc.Login(r.Context(), user.Email, user.Password)
		switch {
		case errors.Is(err, api.ErrUnauthorized) || errors.Is(err, domain.ErrNotFound):
			response.RespondError(r.Context(), w, err, "user not authorised", http.StatusUnauthorized)
			return
		case errors.Is(err, api.ErrTooManyAttempts):
			response.RespondError(r.Context(), w, err, "too many attempts. please reset password", http.StatusTooManyRequests)
			return
		case err != nil:
			response.RespondError(r.Context(), w, err, "user login failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, userRes, http.StatusOK)
	}
}
