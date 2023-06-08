package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"accountID"`
	jwt.RegisteredClaims
}
