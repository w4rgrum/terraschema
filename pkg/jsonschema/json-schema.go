package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/AislingHPE/TerraSchema/pkg/reader"
)

func CreateSchema(path string, strict bool) (string, error) {
	schemaOut := make(map[string]any)

	varMap, err := reader.GetVarMap(path)
	// GetVarMaps returns an error if no .tf files are found in the directory. We
	// ignore this error for now.
	if err != nil && !errors.Is(err, reader.ErrFilesNotFound) {
		return "", fmt.Errorf("error reading tf files at %s: %w", path, err)
	}

	if len(varMap) == 0 {
		return "{}", nil
	}

	schemaOut["$schema"] = "http://json-schema.org/draft-07/schema#"

	if strict {
		schemaOut["additionalProperties"] = false
	}

	properties := make(map[string]any)
	requiredArray := []string{}
	for name, variable := range varMap {
		if variable.Required {
			requiredArray = append(requiredArray, name)
		}
		tc, err := getTypeConstraint(variable.Variable.Type)
		if err != nil {
			return "", fmt.Errorf("getting type constraint for %s: %w", name, err)
		}

		nullableIsTrue := variable.Variable.Nullable != nil && *variable.Variable.Nullable
		node, err := getNodeFromType(name, tc, nullableIsTrue, strict)
		if err != nil {
			return "", fmt.Errorf("%s: %w", name, err)
		}

		def, err := expressionToJSONObject(variable.Variable.Default)
		if err != nil {
			return "", fmt.Errorf("error converting default value to JSON object: %w", err)
		}
		node["default"] = def

		node["description"] = variable.Variable.Description

		properties[name] = node
	}

	schemaOut["properties"] = properties

	slices.Sort(requiredArray) // get required in alphabetical order
	schemaOut["required"] = requiredArray

	out, err := json.MarshalIndent(schemaOut, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(out), nil
}
