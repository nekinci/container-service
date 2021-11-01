package api

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var ErrExpiredToken = errors.New("Token is expired\n")
var ErrInvalidToken = errors.New("Token is invalid\n")

type TokenPayload struct {
	UserId    string `json:"user_id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
}

type TokenResponse struct {
	Token        string    `json:"token"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string    `json:"refresh_token"`
}

func NewTokenResponse(token, refreshToken string, expiresAt time.Time) *TokenResponse {
	return &TokenResponse{
		Token:        token,
		ExpiresAt:    expiresAt,
		RefreshToken: refreshToken,
	}
}

func NewTokenPayload(userId, email string, expiresAt int64) *TokenPayload {
	return &TokenPayload{
		UserId:    userId,
		Email:     email,
		ExpiresAt: expiresAt,
	}
}

func (tokenPayload *TokenPayload) Valid() error {
	t := time.Unix(tokenPayload.ExpiresAt, 0)
	if time.Now().After(t) {
		return ErrExpiredToken
	}

	return nil
}

func generateToken(u User) *TokenResponse {
	var mySignInKey = []byte(os.Getenv("JWT_SECRET"))
	expires := time.Now().Add(time.Hour * 72)
	id := string(u.Id[:])
	tokenPayload := NewTokenPayload(id, u.Email, expires.Unix())

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenPayload)
	token, err := jwtToken.SignedString(mySignInKey)
	if err != nil {
		fmt.Print(err)
		return nil
	}

	refreshToken := generateRefreshToken()

	return NewTokenResponse(token, refreshToken, expires)
}

func generateRefreshToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func validateToken(token string) (*TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &TokenPayload{}, keyFunc)
	if err != nil {
		ver, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(ver.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*TokenPayload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
