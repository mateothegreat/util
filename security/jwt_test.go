package security

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type CustomData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestSignCustomClaims(t *testing.T) {
	GenerateKeySaved()
	privateKey, err := GenerateKey()
	assert.NoError(t, err)

	privateKeyString := fmt.Sprintf("%v", privateKey)
	log.Printf("Generated ECDSA private key: %s", privateKeyString)

	customData := CustomData{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	tokenString, err := SignCustomClaims(customData, time.Hour, privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	publicKey := &privateKey.PublicKey
	claims, err := ParseCustomClaims[CustomData](tokenString, publicKey)
	assert.NoError(t, err)
	assert.Equal(t, customData, claims.Data)
}

func TestParseCustomClaims(t *testing.T) {
	privateKey, err := GenerateKey()
	assert.NoError(t, err)

	customData := CustomData{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	tokenString, err := SignCustomClaims(customData, time.Hour, privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	publicKey := &privateKey.PublicKey
	claims, err := ParseCustomClaims[CustomData](tokenString, publicKey)
	assert.NoError(t, err)
	assert.Equal(t, customData, claims.Data)
}

func TestSign(t *testing.T) {
	privateKey, err := GenerateKey()
	assert.NoError(t, err)

	userId := "1234567890"

	tokenString, err := Sign(userId, time.Hour, privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	t.Logf("Token: %s", tokenString)
	publicKey := &privateKey.PublicKey
	claims, err := Parse(tokenString, publicKey)
	assert.NoError(t, err)
	assert.Equal(t, userId, claims.ID)
}

func TestParse(t *testing.T) {
	privateKey, err := GenerateKey()
	assert.NoError(t, err)

	userId := "1234567890"

	tokenString, err := Sign(userId, time.Hour, privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	publicKey := &privateKey.PublicKey
	claims, err := Parse(tokenString, publicKey)
	assert.NoError(t, err)
	assert.Equal(t, userId, claims.ID)
}

func TestGenerateKey(t *testing.T) {
	privateKey, err := GenerateKey()
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
}

func TestEncodePrivateKey(t *testing.T) {
	privateKey, err := GenerateKey()
	assert.NoError(t, err)

	encoded, err := EncodePrivateKey(privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := DecodePrivateKey(encoded)
	assert.NoError(t, err)
	assert.Equal(t, privateKey, decoded)
}

func TestEncodePublicKey(t *testing.T) {
	privateKey, err := GenerateKey()
	assert.NoError(t, err)

	encoded, err := EncodePublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := DecodePublicKey(encoded)
	assert.NoError(t, err)
	assert.Equal(t, &privateKey.PublicKey, decoded)
}

func TestDecodePrivateKey(t *testing.T) {
	privateKey, err := GenerateKey()
	assert.NoError(t, err)

	encoded, err := EncodePrivateKey(privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := DecodePrivateKey(encoded)
	assert.NoError(t, err)
	assert.Equal(t, privateKey, decoded)
}

func TestDecodePublicKey(t *testing.T) {
	privateKey, err := GenerateKey()
	assert.NoError(t, err)

	encoded, err := EncodePublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := DecodePublicKey(encoded)
	assert.NoError(t, err)
	assert.Equal(t, &privateKey.PublicKey, decoded)
}

func TestInvalidPublicKeyType(t *testing.T) {
	_, err := DecodePublicKey("invalid")
	assert.Error(t, err)
	assert.Equal(t, "illegal base64 data at input byte 4", err.Error())
}
