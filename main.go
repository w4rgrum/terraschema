package main

import (
	"fmt"
	"os"

	"github.com/AislingHPE/TerraSchema/pkg/jsonschema"
)

func main() {
	path := os.Args[1] // absolute path
	strict := false
	if len(os.Args) > 2 && os.Args[2] == "-strict" { // TODO use cobra or flag
		strict = true
	}
	output, err := jsonschema.CreateSchema(path, strict)
	if err != nil {
		fmt.Println(err)

		return
	}
	fmt.Println(output)
}
