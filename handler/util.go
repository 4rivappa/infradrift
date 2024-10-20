package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func CreateFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, PermissionMode)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	return nil
}

func ReadStateFile(stateFile string) (State, error) {
	file, err := os.Open(stateFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return State{}, errors.New("error while opening given state file")
	}
	defer file.Close()

	var state State
	if err := json.NewDecoder(file).Decode(&state); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return State{}, errors.New("error decoding the state file json into structs")
	}
	return state, nil
}

func EraseImportedStateFile() error {
	importedStateFilePath := filepath.Join(FolderName, ImportStateFileName)
	err := os.WriteFile(importedStateFilePath, []byte{}, PermissionMode)
	if err != nil {
		return fmt.Errorf("error erasing contents of imported state file: %w", err)
	}
	return nil
}

func GetProviderTF(importStateFilePath string) []byte {
	providerContent := `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.63.0"
    }
  }
  backend "local" {
    path = "%v"
  }
}

provider "aws" {
  region = "us-east-1"
}`

	return []byte(fmt.Sprintf(providerContent, importStateFilePath))
}

func stringInSliceOfStrings(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
