package config

import "os"

// JWTSecret returns the single source of truth for the JWT signing key.
func JWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "krishisetu_secret_key" // fallback for dev
	}
	return []byte(secret)
}
