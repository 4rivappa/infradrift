package handler

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

var FolderName string = "import-infra"
var PermissionMode fs.FileMode = 0755

var ProviderFileName string = "provider.tf"
var MainFileName string = "main.tf"
var ImportStateFileName string = "terraform.tfstate"

var unsupportedResources = []string{"aws_caller_identity", "aws_partition", "aws_region", "aws_iam_policy_document"}
var resourceIdOutliners = map[string]RequiredIdsResource{
	"aws_iam_role_policy_attachment": {
		MainFileIds: []string{"role", "policy_arn"},
		CommandId:   "policy_arn",
	},
}

func setupImportRequirements() error {
	for _, folderName := range []string{"", "/diff/import", "/diff/state"} {
		err := CreateFolder(FolderName + folderName)
		if err != nil {
			return err
		}
	}

	providerFilePath := filepath.Join(FolderName, ProviderFileName)
	fmt.Println(providerFilePath)
	err := os.WriteFile(providerFilePath, GetProviderTF(ImportStateFileName), PermissionMode)
	if err != nil {
		return fmt.Errorf("error creating provider.tf: %w", err)
	}

	mainFilePath := filepath.Join(FolderName, MainFileName)
	fmt.Println(mainFilePath)
	err = os.WriteFile(mainFilePath, []byte{}, PermissionMode)
	if err != nil {
		return fmt.Errorf("error creating main.tf: %w", err)
	}

	cmd := exec.Command("terraform", "init")
	cmd.Dir = FolderName
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing terraform init: %w, output: %s", err, output)
	}

	return nil
}

func HandleDrift(stateFile string) error {

	state, err := ReadStateFile(stateFile)
	if err != nil {
		return err
	}
	// fmt.Println(state)
	fmt.Println("===================================")

	setupErr := setupImportRequirements()
	if setupErr != nil {
		fmt.Println(setupErr)
	}

	indexCounter := 0

	for _, resource := range state.Resources {

		if stringInSliceOfStrings(resource.Type, unsupportedResources) {
			continue
		}

		for _, instance := range resource.Instances {

			loadInstanceErr := LoadInstanceHandler(resource.Type, resource.Name, instance)
			if loadInstanceErr != nil {
				fmt.Printf("Problem in loading an instance: %v", loadInstanceErr)
			}

			indexCounter += 1
			instanceDiff, compareImportedErr := compareImportedHandler(resource.Type, resource.Name, instance, indexCounter)
			if compareImportedErr != nil {
				fmt.Printf("Problem in comparing imported instance: %v", compareImportedErr)
			}

			fmt.Printf("%d | %s | %s | %v\n", indexCounter, resource.Type, resource.Name, instanceDiff)

			err := EraseImportedStateFile()
			if err != nil {
				fmt.Printf("Error - while erasing imported state file")
			}

			// updated the above to fetch difference in keys of instances
			// print them along with num of that instance, processed

			// empty terraform.tfstate file, before re-executing for another instance
		}
	}

	return nil
}
