package challenge

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math"
	"strconv"
)

// Generate returns a random base64 string of length n.
func Generate(n int) string {
	// one byte is 1/3 wider than a base 64 value, 2^8 vs 2^6.
	buff := make([]byte, int(math.Ceil(float64(n)/1.33333)))
	_, _ = rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:n]
}

func zeros(prefix string) bool {
	for _, ch := range prefix {
		if ch != '0' {
			return false
		}
	}
	return true
}

func Solve(data string, bits int) string {
	var solution int
	for {
		token := data + strconv.Itoa(solution)
		hash := sha256.Sum256([]byte(token))
		hashHex := hex.EncodeToString(hash[:])
		if zeros(hashHex[:bits]) {
			return strconv.Itoa(solution)
		}
		solution++
	}
}

func Verify(data, solution string, bits int) bool {
	token := data + solution
	hash := sha256.Sum256([]byte(token))
	hashHex := hex.EncodeToString(hash[:])
	return zeros(hashHex[:bits])
}
