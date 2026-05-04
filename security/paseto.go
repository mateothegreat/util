package security

import (
	"time"

	"aidanwoods.dev/go-paseto"
)

// PasetoSign signs a token with the given key and data.
//
// Arguments:
//   - key: The secret key to use for signing the token.
//   - data: The data to sign.
//   - exp: The expiration time of the token.
//
// Returns:
//   - The signed token string.
func PasetoSign(key paseto.V4AsymmetricSecretKey, data interface{}, exp time.Time) string {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.Set("data", data)
	token.SetExpiration(exp)

	return token.V4Sign(key, nil)
}

// PasetoParse parses a token and returns the data or an error.
//
// Arguments:
//   - key: The secret key to use for parsing the token.
//   - token: The token string to parse.
//
// Returns:
//   - The parsed data.
//   - An error if the token is invalid.
func PasetoParse[T any](key paseto.V4AsymmetricSecretKey, token string) (T, error) {
	var t T

	publicKey := key.Public() // Extract the public key from the secret key
	parser := paseto.NewParser()
	parsed, err := parser.ParseV4Public(publicKey, token, nil)
	if err != nil {
		return t, err
	}

	err = parsed.Get("data", &t)
	if err != nil {
		return t, err
	}

	return t, nil
}
