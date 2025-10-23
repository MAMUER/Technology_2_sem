package jwt

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type Validator interface {
	SignAccess(userID int64, email, role string) (string, error)
	SignRefresh(userID int64) (string, error)
	Parse(tokenStr string) (jwt.MapClaims, error)
}

type HS256 struct {
	secret []byte
	accessTTL time.Duration
	refreshTTL time.Duration
}

func NewHS256(secret []byte, accessTTL, refreshTTL time.Duration) *HS256 {
	return &HS256{
		secret: secret, 
		accessTTL: accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (h *HS256) SignAccess(userID int64, email, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"type":  "access",
		"iat":   now.Unix(),
		"exp":   now.Add(h.accessTTL).Unix(),
		"iss":   "pz10-auth",
		"aud":   "pz10-clients",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(h.secret)
}

func (h *HS256) SignRefresh(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  userID,
		"type": "refresh",
		"iat":  now.Unix(),
		"exp":  now.Add(h.refreshTTL).Unix(),
		"iss":  "pz10-auth",
		"aud":  "pz10-clients",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(h.secret)
}

func (h *HS256) Parse(tokenStr string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return h.secret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !t.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	
	return claims, nil
}