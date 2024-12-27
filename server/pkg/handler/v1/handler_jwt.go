package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/neuro-soup/evochi/server/internal/worker"
)

type Claims struct {
	jwt.RegisteredClaims
}

func (c *Claims) Admin() bool {
	return c.Subject == "admin"
}

// jwtKey returns the key used to verify the JWT token.
func (h *Handler) jwtKey(t *jwt.Token) (any, error) {
	return []byte(h.config.JWTSecret), nil
}

// authenticate reads the JWT token from the "Authorization" header (Bearer)
// and verifies it.
func (h *Handler) authenticate(header http.Header) (*Claims, error) {
	tok := strings.TrimPrefix(header.Get("Authorization"), "Bearer ")
	if tok == "" {
		return nil, errors.New("auth: token is empty")
	}

	token, err := jwt.ParseWithClaims(tok, &Claims{}, h.jwtKey)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("auth: token is invalid")
	}

	return token.Claims.(*Claims), nil
}

// authenticateWorker reads the JWT token from the "Authorization" header (Bearer)
// and verifies it. It also checks that the worker exists and, if so, returns
// the worker and the claims.
func (h *Handler) authenticateWorker(
	header http.Header,
) (*worker.Worker, *Claims, error) {
	claims, err := h.authenticate(header)
	if err != nil {
		return nil, nil, err
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, nil, fmt.Errorf("auth: failed to parse worker ID: %w", err)
	}

	w := h.workers.Get(id)
	if w == nil {
		return nil, nil, errors.New("auth: worker not found")
	}

	return w, claims, nil
}

// // createJWT creates a new JWT token for the worker.
func (h *Handler) createJWT(w *worker.Worker) (*jwt.Token, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "evochi@v1",
			Subject:   w.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 7)), // TODO: make configurable
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodPS512.SigningMethodRSA, claims), nil // TODO: make method configurable
}
