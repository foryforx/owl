package httpapi

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/foryforx/owl/owl-auth/internal/domain"
	"github.com/foryforx/owl/owl-auth/internal/transport/httpapi/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Middleware func(Handler) Handler

func WithLogging() Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debugln("api called", "url", r.Method, r.RequestURI)
			traceID := uuid.New()
			ctx := context.WithValue(r.Context(), domain.CKey("traceID"), traceID)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
			log.Debugln("api finished", "url", r.Method, r.RequestURI)
		})
	}
}

func ChainMiddleware(mm ...Middleware) Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, m := range mm {
				h = m(h)
			}
			h.ServeHTTP(w, r)
		})
	}
}

func WithAuthenticator(usrSvc usersGetter) Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := GetTokenFromRequest(r)

			if err != nil {
				response.RespondError(r.Context(), w, err, "Invalid token", http.StatusUnauthorized)
				return
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			token, err := parseToken(tokenString)
			if err != nil {
				response.RespondError(r.Context(), w, err, "Invalid token", http.StatusUnauthorized)
				return
			}
			claims := token.Claims.(*domain.Claims)
			user, err := usrSvc.GetUser(r.Context(), claims.ID, claims.AccountID)
			if err != nil {
				response.RespondError(r.Context(), w, err, "Invalid token", http.StatusUnauthorized)
				return
			}
			userD := &domain.User{
				ID:        user.ID,
				AccountID: user.AccountID,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}
			ctx := context.WithValue(r.Context(), domain.CKey("user"), userD)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func GetTokenFromRequest(r *http.Request) (token string, err error) {
	token = r.Header.Get("Authorization")
	if len(token) > 0 {
		return token, nil
	}
	err = errors.New("JWT token is not found on the request")
	return "", err
}
func parseToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := errors.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, err
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
