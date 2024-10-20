package handler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func genMainFileContent(resourceType, resourceName string, ids []string, instance Instance) (string, error) {
	placeHolderContent := fmt.Sprintf(`resource "%s" "%s" {`, resourceType, resourceName)

	for _, id := range ids {
		value, exists := instance.Attributes[id]
		if !exists {
			return "", errors.New("no id found in instance attributes")
		}

		idValue, ok := value.(string)
		if !ok {
			return "", errors.New("id attribute's value is not string - assertion error")
		} else {
			placeHolderContent += "\n" + fmt.Sprintf(`  %s = "%s"`, id, idValue)
		}
	}
	placeHolderContent += "\n}"
	return placeHolderContent, nil
}

func writeMainFile(data []byte) error {
	mainFilePath := filepath.Join(FolderName, MainFileName)
	err := os.WriteFile(mainFilePath, data, PermissionMode)
	if err != nil {
		return fmt.Errorf("error writing to main.tf: %w", err)
	}
	return nil
}

func LoadInstanceHandler(resourceType string, resourceName string, instance Instance) error {
	resourceId := []string{"id"}
	if id, exists := resourceIdOutliners[resourceType]; exists {
		resourceId = id.MainFileIds
	}

	value, exists := instance.Attributes[resourceId[0]]
	if !exists {
		return errors.New("no id found in instance attributes")
	}

	idValue, ok := value.(string)
	if !ok {
		return errors.New("id attribute's value is not string - assertion error")
	}
	mainFileContent, genErr := genMainFileContent(resourceType, resourceName, resourceId, instance)
	if genErr != nil {
		return genErr
	}

	writeErr := writeMainFile([]byte(mainFileContent))
	if writeErr != nil {
		return writeErr
	}

	if id, exists := resourceIdOutliners[resourceType]; exists {
		requiredIdForCmd := id.CommandId
		value, exists := instance.Attributes[requiredIdForCmd]
		newIdValue, ok := value.(string)
		if !ok {
			return errors.New("new id attribute's value is not string - assertion error")
		}
		if exists {
			if resourceType == "aws_iam_role_policy_attachment" {
				idValue += "/" + newIdValue
			} else {
				idValue = newIdValue
			}
		}
	}
	cmd := exec.Command("terraform", "import", fmt.Sprintf("%s.%s", resourceType, resourceName), idValue)
	cmd.Dir = FolderName
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing terraform import: %w, output: %s", err, output)
	}

	return nil
}
