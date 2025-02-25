package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

var signification string

func hashIDAndEmail(id uint, email string) string {
	idStr := strconv.Itoa(int(id))
	combined := idStr + email
	hasher := sha256.New()
	hasher.Write([]byte(combined))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	signification = hashString
	return hashString
}
