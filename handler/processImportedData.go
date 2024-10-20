package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

func filterStateForReqInstance(state State, resourceType string, resourceName string) (Instance, error) {
	for _, resource := range state.Resources {
		if resource.Type == resourceType && resource.Name == resourceName {
			if len(resource.Instances) > 0 {
				return resource.Instances[0], nil
			}
		}
	}
	return Instance{}, errors.New("couldn't find required instance from imported state file")
}

func backupInstances(inst1, inst2 Instance, file1, file2 string) error {
	data1, err := json.MarshalIndent(inst1, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling instance 1: %v", err)
	}
	if err := os.WriteFile(file1, data1, PermissionMode); err != nil {
		return fmt.Errorf("error writing to file %s: %v", file1, err)
	}

	data2, err := json.MarshalIndent(inst2, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling instance 2: %v", err)
	}
	if err := os.WriteFile(file2, data2, PermissionMode); err != nil {
		return fmt.Errorf("error writing to file %s: %v", file2, err)
	}

	return nil
}

func diffInstances(first, second Instance) []string {
	mismatchedKeys := []string{}
	for key, value := range first.Attributes {
		if val2, ok := second.Attributes[key]; ok {
			if !reflect.DeepEqual(value, val2) {
				mismatchedKeys = append(mismatchedKeys, key)
			}
		}
	}
	for key := range second.Attributes {
		if _, ok := first.Attributes[key]; !ok {
			mismatchedKeys = append(mismatchedKeys, key)
		}
	}
	return mismatchedKeys
}

func compareImportedHandler(resourceType, resourceName string, instance Instance, indexCounter int) ([]string, error) {
	importedStateFilePath := filepath.Join(FolderName, ImportStateFileName)
	state, readStateFileErr := ReadStateFile(importedStateFilePath)
	if readStateFileErr != nil {
		return []string{}, readStateFileErr
	}

	importedInstance, filterErr := filterStateForReqInstance(state, resourceType, resourceName)
	if filterErr != nil {
		return []string{}, filterErr
	}

	instanceFileName := fmt.Sprintf(`%d.json`, indexCounter)
	stateInstanceFilePath := filepath.Join(FolderName, "diff", "state", instanceFileName)
	importedInstanceFilePath := filepath.Join(FolderName, "diff", "import", instanceFileName)
	backupErr := backupInstances(instance, importedInstance, stateInstanceFilePath, importedInstanceFilePath)
	if backupErr != nil {
		return []string{}, backupErr
	}

	diff := diffInstances(instance, importedInstance)
	return diff, nil
}
