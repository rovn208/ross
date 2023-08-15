package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const minSecretKeySize = 10

type JWTMaker struct {
	secretKey string
}

type JWTCustomClaims struct {
	Username string `json:"username,omitempty"`
	UserID   int64  `json:"user_id,omitempty"`
	jwt.RegisteredClaims
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and userID
func (maker *JWTMaker) CreateToken(username string, userID int64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, userID, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, convertPayloadToClaims(payload))
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// VerifyToken verifies the token Payload if it is valid
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JWTCustomClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return convertClaimsToPayload(claims), nil
}

func convertPayloadToClaims(payload *Payload) JWTCustomClaims {
	return JWTCustomClaims{
		payload.Username,
		payload.UserID,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			ID:        payload.ID.String(),
		}}
}

func convertClaimsToPayload(claims *JWTCustomClaims) *Payload {
	return &Payload{
		Username:  claims.Username,
		UserID:    claims.UserID,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiredAt: claims.ExpiresAt.Time,
	}
}
