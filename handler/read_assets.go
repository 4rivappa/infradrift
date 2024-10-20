package handler

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadAsset(filename string) ([]byte, error) {
	assetPath := filepath.Join("assets", filename)
	data, err := os.ReadFile(assetPath)
	if err != nil {
		return nil, fmt.Errorf("error reading asset file: %w", err)
	}
	return data, nil
}
