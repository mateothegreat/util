package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

// GenerateKey generates a new ECDSA key pair.
func GenerateKey() (*ecdsa.PrivateKey, error) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func GenerateKeySaved() (*ecdsa.PrivateKey, error) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	// Encode private key to base64
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	privateKeyBase64 := base64.URLEncoding.EncodeToString(privateKeyBytes)

	log.Printf("Generated ECDSA private key: %s", privateKeyBase64)
	// Encode public key to base64
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	publicKeyBase64 := base64.URLEncoding.EncodeToString(publicKeyBytes)

	log.Printf("Generated ECDSA public key: %s", publicKeyBase64)

	// Save private key to file
	err = ioutil.WriteFile("private.key", []byte(privateKeyBase64), 0600)
	if err != nil {
		return nil, err
	}

	// Save public key to file
	err = ioutil.WriteFile("public.key", []byte(publicKeyBase64), 0644)
	if err != nil {
		return nil, err
	}

	fmt.Println("Keys saved to files.")
	return privateKey, nil

}

// EncodePrivateKey encodes the private key to a base64 string.
func EncodePrivateKey(privateKey *ecdsa.PrivateKey) (string, error) {
	der, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(der), nil
}

// DecodePrivateKey decodes the private key from a base64 string.
func DecodePrivateKey(encoded string) (*ecdsa.PrivateKey, error) {
	der, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParseECPrivateKey(der)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// EncodePublicKey encodes the public key to a base64 string.
func EncodePublicKey(publicKey *ecdsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(der), nil
}

// DecodePublicKey decodes the public key from a base64 string.
func DecodePublicKey(encoded string) (*ecdsa.PublicKey, error) {
	der, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	publicKey, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	switch publicKey := publicKey.(type) {
	case *ecdsa.PublicKey:
		return publicKey, nil
	default:
		return nil, errors.New("invalid public key type")
	}
}
