package testsetup

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv attempts to load a .env file from the current working directory
// or from common parent paths (package test execution changes working dir).
// It returns the path that was loaded or an empty string if none found.
func LoadEnv() string {
	// Try common locations relative to package directories when running `go test ./...`
	candidates := []string{
		".env",
		"../../.env",
		"../.env",
		"../../../.env",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err == nil {
				return p
			}
		}
	}
	// Try automatic load (searches up to root)
	if err := godotenv.Load(); err == nil {
		return ".env"
	}
	fmt.Println("testsetup: .env not found in candidates; proceeding without loading")
	return ""
}
