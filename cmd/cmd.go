// (C) Copyright 2024 Hewlett Packard Enterprise Development LP
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/HewlettPackard/terraschema/pkg/jsonschema"
)

// wanted behaviour:
// - disallow-additional-properties: disallow additional properties in schema (default is false)
// - overwrite: overwrite an existing file (default is false for safety reasons)
// - stdout: suppress errors and output schema to stdout (generally not recommended)
// - output: file, default is ./schema.json. Allow creation of directories.
// - input: folder, default is .
// - allow-empty: if no variables are found, print empty schema and exit with 0

var (
	disallowAdditionalProperties bool
	overwrite                    bool
	allowEmpty                   bool
	outputStdOut                 bool
	output                       string
	input                        string
)

// rootCmd is the base command for terraschema
var rootCmd = &cobra.Command{
	Use:   "terraschema",
	Short: "Generate JSON schema from HCL Variable Blocks in a Terraform/OpenTofu module",
	Long:  `TODO: Long description`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := filepath.Abs(input) // absolute path
		if err != nil {
			fmt.Printf("could not get absolute path: %v\n", err)
			os.Exit(1)
		}
		output, err := jsonschema.CreateSchema(path, false)
		if err != nil {
			fmt.Printf("error creating schema: %v\n", err)
			os.Exit(1)
		}
		jsonOutput, err := json.MarshalIndent(output, "", "    ")
		if err != nil {
			fmt.Printf("error marshalling schema: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonOutput))
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// TODO: implement
	rootCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing schema file")
	// TODO: implement
	rootCmd.Flags().BoolVar(&outputStdOut, "stdout", false,
		"output schema content to stdout instead of a file and disable error output",
	)
	// TODO: implement
	rootCmd.Flags().BoolVar(&disallowAdditionalProperties, "disallow-additional-properties", false,
		"set additionalProperties to false in the generated schema and in nested objects",
	)
	// TODO: implement
	rootCmd.Flags().BoolVar(&allowEmpty, "allow-empty", false, "allow empty schema if no variables are found, otherwise error")
	rootCmd.Flags().StringVarP(&input, "input", "i", ".", "input folder containing .tf files")
	// TODO: implement
	rootCmd.Flags().StringVarP(&output, "output", "o", "schema.json", "output file path for schema")
}
