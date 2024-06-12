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
	output := jsonschema.CreateSchema(path, strict)
	fmt.Println(output)
}
