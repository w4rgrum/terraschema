package main

import (
	"fmt"
	"os"

	"github.com/AislingHPE/terraschema/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("exited with error: %v\n", err)
		os.Exit(1)
	}
}
