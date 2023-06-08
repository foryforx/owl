package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/foryforx/owl/owl-auth/api"
	"github.com/foryforx/owl/owl-auth/internal/domain"
	"github.com/foryforx/owl/owl-auth/internal/transport/httpapi/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AccountRouter struct {
	AccountService api.AccountService
}

func (a *AccountRouter) Routes() []Route {
	return []Route{
		{http.MethodPost, "/accounts", CreateAccountHandler(a.AccountService), false},
		{http.MethodGet, "/accounts/{id}", GetAccountHandler(a.AccountService), false},
		{http.MethodGet, "/accounts", GetAccountsHandler(a.AccountService), false},
		{http.MethodPut, "/accounts/{id}", UpdateAccountHandler(a.AccountService), false},
		{http.MethodDelete, "/accounts/{id}", DeleteAccountHandler(a.AccountService), false},
	}
}

type accountsGetter interface {
	GetAccount(ctx context.Context, id uuid.UUID) (*api.AccountResponse, error)
	GetAccounts(ctx context.Context) ([]*api.AccountResponse, error)
}

type accountsCreator interface {
	CreateAccount(ctx context.Context, account *api.AccountRequest) (uuid.UUID, error)
}

type accountsUpdator interface {
	UpdateAccount(ctx context.Context, account *api.AccountRequest) error
}

type accountsDeleter interface {
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

func CreateAccountHandler(accSvc accountsCreator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var account api.AccountRequest

		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			response.RespondError(r.Context(), w, err, "invalid request", http.StatusBadRequest)
			return
		}

		id, err := accSvc.CreateAccount(r.Context(), &account)

		switch {
		case errors.As(err, &domain.DuplicateEntryError{}):
			response.RespondError(r.Context(), w, err, "name already exists", http.StatusConflict)
			return
		case err != nil:
			response.RespondError(r.Context(), w, err, "registration failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, api.CreateAccountResponse{ID: id}, http.StatusCreated)
	}
}

func GetAccountHandler(accSvc accountsGetter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := r.Context().Value(domain.CKey("user")).(*domain.User)
		if !ok {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		idStr := mux.Vars(r)["id"]
		if idStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "account doesnt exists", http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "account doesnt exists", http.StatusNotFound)
			return
		}
		if currentUser.AccountID != id {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		account, err := accSvc.GetAccount(r.Context(), id)
		switch {
		case errors.Is(err, domain.ErrNotFound):
			response.RespondError(r.Context(), w, err, "account not found", http.StatusNotFound)
			return
		case err != nil:
			response.RespondError(r.Context(), w, err, "get accounts failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, account, http.StatusOK)
	}
}

func GetAccountsHandler(accSvc accountsGetter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := accSvc.GetAccounts(r.Context())
		switch {
		case errors.Is(err, domain.ErrNotFound):
			response.RespondError(r.Context(), w, err, "account not found", http.StatusNotFound)
			return
		case err != nil:
			response.RespondError(r.Context(), w, err, "get accounts failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, accounts, http.StatusOK)
	}
}

func UpdateAccountHandler(accSvc accountsUpdator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := r.Context().Value(domain.CKey("user")).(*domain.User)
		if !ok {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		var account api.AccountRequest
		idStr := mux.Vars(r)["id"]
		if idStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "account doesnt exists", http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "account doesnt exists", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			response.RespondError(r.Context(), w, err, "invalid request", http.StatusBadRequest)
			return
		}
		account.ID = id
		if currentUser.AccountID != id {
			response.RespondError(r.Context(), w, errors.New("unauthorized access"), "unauthorized access", http.StatusUnauthorized)
			return
		}
		if err := accSvc.UpdateAccount(r.Context(), &account); err != nil {
			response.RespondError(r.Context(), w, err, "update failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, nil, http.StatusOK)
	}
}

func DeleteAccountHandler(accSvc accountsDeleter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		if idStr == "" {
			response.RespondError(r.Context(), w, errors.New("invalid request"), "account doesnt exists", http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.RespondError(r.Context(), w, err, "account doesnt exists", http.StatusNotFound)
			return
		}
		if err := accSvc.DeleteAccount(r.Context(), id); err != nil {
			response.RespondError(r.Context(), w, err, "delete failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, nil, http.StatusOK)
	}
}
