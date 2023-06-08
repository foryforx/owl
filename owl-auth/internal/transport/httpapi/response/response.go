package response

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/foryforx/owl/owl-auth/internal/domain"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func RespondError(ctx context.Context, w http.ResponseWriter, err error, friendlyMessage string, status int) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	traceID, _ := ctx.Value(domain.CKey("traceID")).(uuid.UUID)
	log.Errorln(traceID, err)
	_ = json.NewEncoder(w).Encode(
		Error{
			Code:    status,
			Message: friendlyMessage,
			Err:     err,
		},
	)
}

func RespondInternalServerError(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	traceID, _ := ctx.Value(domain.CKey("traceID")).(uuid.UUID)
	log.Errorln(traceID, err)
	_ = json.NewEncoder(w).Encode(
		Error{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong. Please try again later.",
			Err:     err,
		},
	)
}
