package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func GenerateOTP() (string, time.Time) {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	otp := n.Int64() + 100000 // 6-digit
	expiry := time.Now().Add(10 * time.Minute)
	return fmt.Sprintf("%d", otp), expiry
}
