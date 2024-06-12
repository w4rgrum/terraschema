package reader

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/AislingHPE/TerraSchema/pkg/model"
)

var fileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
	},
}

var ErrFilesNotFound = fmt.Errorf("no .tf files found in directory")

func GetVarMap(path string) (map[string]model.TranslatedVariable, error) {
	// read all tf files in directory
	files, err := filepath.Glob(filepath.Join(path, "*.tf"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, ErrFilesNotFound
	}

	parser := hclparse.NewParser()

	varMap := make(map[string]model.TranslatedVariable)
	for _, fileName := range files {
		file, d := parser.ParseHCLFile(fileName)
		if d.HasErrors() {
			return nil, d
		}

		blocks, _, d := file.Body.PartialContent(fileSchema)
		if d.HasErrors() {
			return nil, d
		}
		for _, block := range blocks.Blocks {
			name, translated, err := getTranslatedVariableFromBlock(block, file)
			if err != nil {
				return nil, err
			}
			varMap[name] = translated
		}
	}

	return varMap, nil
}

func getTranslatedVariableFromBlock(block *hcl.Block, file *hcl.File) (string, model.TranslatedVariable, error) {
	attributes, d := block.Body.JustAttributes()
	if d.HasErrors() {
		return "", model.TranslatedVariable{}, d
	}

	name := block.Labels[0]
	variable := model.Variable{}
	d = gohcl.DecodeBody(block.Body, nil, &variable)
	if d.HasErrors() {
		return name, model.TranslatedVariable{}, d
	}

	out := model.TranslatedVariable{Variable: variable, Required: true}

	// Get type, default, and condition as strings and add them to the translated variable struct.
	// This is to make the code easier to debug, since hcl.Expressions are difficult to read out of context.
	// Also filter any values for hcl.Expressions whose corresponding fields are not present in the block, and set them
	// to nil.

	// check if 'default' exists in the block directly
	if _, ok := attributes["default"]; ok {
		out.DefaultAsString = expressionAsStringPointer(variable.Default, file)
		out.Required = false
	} else {
		// do not keep a hcl expression for 'default' if no 'default' field is present
		out.Variable.Default = nil
	}

	// check if 'type' exists in the block directly
	if _, ok := attributes["type"]; ok {
		out.TypeAsString = expressionAsStringPointer(variable.Type, file)
	} else {
		// do not keep a hcl expression for 'type' if no 'type' field is present
		out.Variable.Type = nil
	}

	// if a validation block does not exist, variable.Validation is nil.
	if variable.Validation != nil && variable.Validation.Condition != nil {
		out.ConditionAsString = expressionAsStringPointer(variable.Validation.Condition, file)
	}

	return name, out, nil
}

func expressionAsStringPointer(in hcl.Expression, f *hcl.File) *string {
	out := string(in.Range().SliceBytes(f.Bytes))

	return &out
}
