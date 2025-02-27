package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func HashIDAndEmail(id uint, email string) string {
	combined := strconv.Itoa(int(id)) + email
	hasher := sha256.New()
	hasher.Write([]byte(combined))
	return hex.EncodeToString(hasher.Sum(nil))
}
