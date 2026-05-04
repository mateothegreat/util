package security

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomString generates a random string of length n.
// It uses a combination of numbers and uppercase and lowercase letters.
// The generated string is returned along with any error encountered during the process.
func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

func GenerateRandomNumber(min int, max int) int {
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(randomNumber.Int64()) + min
}
