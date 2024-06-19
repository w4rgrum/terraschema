// (C) Copyright 2024 Hewlett Packard Enterprise Development LP
package main

import (
	"fmt"
	"os"

	"github.com/HewlettPackard/terraschema/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("exited with error: %v\n", err)
		os.Exit(1)
	}
}
