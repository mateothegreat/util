package random

import (
	"bytes"
	"crypto/rand"
	"io"
	"math/big"
)

// RandomGenerator generates random data.
type RandomGenerator struct {
	bytes int64
}

// NewRandomGenerator returns a new RandomGenerator.
func NewRandomGenerator(bytes int64) *RandomGenerator {
	return &RandomGenerator{
		bytes: bytes,
	}
}

func (g *RandomGenerator) Reader() (io.Reader, error) {
	data := make([]byte, g.bytes)
	if _, err := rand.Read(data); err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

func GenerateRandomNumber(min int, max int) int {
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(randomNumber.Int64()) + min
}
