package jsonschema

import (
	"github.com/hashicorp/hcl/v2"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// expressionToJSONObject converts an HCL expression to an interface{} that can be marshaled to JSON.
func expressionToJSONObject(in hcl.Expression) (any, error) {
	v, d := in.Value(&hcl.EvalContext{})
	if d.HasErrors() {
		return nil, d
	}

	return ctyjson.SimpleJSONValue{Value: v}, nil
}
