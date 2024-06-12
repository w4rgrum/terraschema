package jsonschema

import (
	"encoding/json"
	"slices"

	"github.com/AislingHPE/TerraSchema/pkg/model"
)

func CreateSchema(path string, strict bool) string {
	schemaOut := make(map[string]any)

	varMap := make(map[string]model.TranslatedVariable) // getVarMap(path)

	if len(varMap) == 0 {
		return "{}"
	}

	schemaOut["$schema"] = "http://json-schema.org/draft-07/schema#"

	if strict {
		schemaOut["additionalProperties"] = false
	}

	properties := make(map[string]any)
	requiredArray := []string{}

	schemaOut["properties"] = properties

	slices.Sort(requiredArray) // get required in alphabetical order
	schemaOut["required"] = requiredArray

	out, err := json.MarshalIndent(schemaOut, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(out)
}
