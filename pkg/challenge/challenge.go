package challenge

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func validPrefix(prefix string) bool {
	for _, ch := range prefix {
		if ch != '0' {
			return false
		}
	}
	return true
}

// Generate returns a random base64 string of length n,
// prefixed by challenge difficulty with ":" delimiter.
// Example: Generate(12, 4) == "4:btY5IVj_8FN7".
func Generate(length, difficulty int) string {
	// one byte is 1/3 wider than a base 64 value, 2^8 vs 2^6.
	buff := make([]byte, int(math.Ceil(float64(length)/1.33333)))
	_, _ = rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return strconv.Itoa(difficulty) + ":" + str[:length]
}

// Solve finds a solution which when appended to data payload
// and hashed, generates a token prefixed with a number of
// leading zeros specified in the challenge prefix.
func Solve(data string) (string, error) {
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return "", errors.New("invalid challenge format")
	}

	difficulty, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("failed to parse difficulty: %w", err)
	}

	var solution int
	for {
		hash := sha256.Sum256([]byte(parts[1] + strconv.Itoa(solution)))
		hashHex := hex.EncodeToString(hash[:])
		if validPrefix(hashHex[:difficulty]) {
			return strconv.Itoa(solution), nil
		}
		solution++
	}
}

// Verify checks that the provided solution
// satisfies the difficulty constraint.
func Verify(data, solution string) error {
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return errors.New("invalid challenge format")
	}

	difficulty, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("failed to parse difficulty: %w", err)
	}

	hash := sha256.Sum256([]byte(parts[1] + solution))
	hashHex := hex.EncodeToString(hash[:])
	if !validPrefix(hashHex[:difficulty]) {
		return fmt.Errorf("incorrect solution %q for challange %q", solution, data)
	}

	return nil
}
