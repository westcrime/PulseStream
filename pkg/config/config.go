package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadDotEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		path := filepath.Join(dir, ".env")
		if _, err = os.Stat(path); err == nil {
			return godotenv.Load(path)
		}
		if _, err = os.Stat(filepath.Join(dir, ".git")); err == nil {
			return nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil
		}
		dir = parent
	}
}
