package main

import (
	"fmt"
	"os"

	"infradrift/handler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: infra-drift <state-file>")
		return
	}

	folderName := os.Args[1]

	if err := handler.HandleDrift(folderName); err != nil {
		fmt.Println(err)
	}
}
