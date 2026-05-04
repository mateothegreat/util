package security

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}
type UserClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func SignCustomClaims[T any](custom T, duration time.Duration, privateKey *ecdsa.PrivateKey) (string, error) {
	claims := CustomClaims[T]{
		custom,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	t, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

// Parse parses the token (validates) and rpeturns the claims.
func ParseCustomClaims[T any](tokenString string, publicKey *ecdsa.PublicKey) (*CustomClaims[T], error) {
	claims := &CustomClaims[T]{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func Sign(id string, duration time.Duration, privateKey *ecdsa.PrivateKey) (string, error) {
	claims := UserClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

// Parse parses the token (validates) and returns the claims.
func Parse(tokenString string, publicKey *ecdsa.PublicKey) (*UserClaims, error) {
	claims := &UserClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
