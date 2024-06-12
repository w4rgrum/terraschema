package jsonschema

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
)

var simpleTypeMap = map[string]string{
	"string": "string",
	"number": "number",
	"bool":   "boolean",
}

// converts the expression into a type constraint and marshals it to JSON.
// Returns the JSON output as a []byte.
// More info on exactly how this works is here:
// https://pkg.go.dev/github.com/zclconf/go-cty@v1.14.4/cty#Type.MarshalJSON
func getTypeConstraint(in hcl.Expression) (any, error) {
	if in == nil {
		return "any", nil
	}

	t, d := typeexpr.TypeConstraint(in)
	if d.HasErrors() {
		return nil, fmt.Errorf("could not parse type constraint from expression: %w", d)
	}
	typeJSON, err := t.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not marshal constraint to JSON: %w", err)
	}

	var typeInterface any
	err = json.Unmarshal(typeJSON, &typeInterface)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal type from JSON: %w", err)
	}

	return typeInterface, nil
}

func getNodeFromType(name string, typeInterface any, nullable bool, strict bool) (map[string]any, error) {
	if nullable {
		return getNullableNode(name, typeInterface, strict)
	}

	switch t := typeInterface.(type) {
	case string:
		if simpleType, ok := simpleTypeMap[t]; ok {
			return map[string]any{"type": simpleType}, nil
		} else {
			return nil, fmt.Errorf("unsupported type %q", t)
		}
	case []any:
		return getNodeFromSlice(t, strict)
	default:
		return nil, fmt.Errorf("unsupported type for %#v", typeInterface)
	}
}

func getNullableNode(name string, typeInterface any, strict bool) (map[string]any, error) {
	node := make(map[string]any)
	if typeInterface == nil {
		return node, nil
	}
	internalNode, err := getNodeFromType(name, typeInterface, false, strict)
	if err != nil {
		return nil, err
	}
	title, ok := internalNode["type"].(string)
	if !ok {
		return nil, fmt.Errorf("could not get type %v as a string", internalNode["type"])
	}

	internalNode["title"] = title

	node["anyOf"] = []any{
		map[string]any{"type": "null", "title": "null"},
		internalNode,
	}
	node["title"] = fmt.Sprintf("%s: Select a type", name)

	return node, nil
}

func getNodeFromSlice(in []any, strict bool) (map[string]any, error) {
	switch in[0] {
	// "object" affects additionalProperties, properties, type and required
	case "object":
		return getObject(in, strict)
	// "map" affects additionalProperties and type.
	case "map":
		return getMap(in, strict)
	// "list" affects items, type
	case "list":
		return getList(in, strict)
	// "set" affects items, type, uniqueItems
	case "set":
		return getSet(in, strict)
	// "tuple" affects items, type, maxItems, minItems
	case "tuple":
		return getTuple(in, strict)
	default:
		panic("unknown type")
	}
}

func getObject(in []any, strict bool) (map[string]any, error) {
	node := map[string]any{
		"type": "object",
	}
	if strict {
		node["additionalProperties"] = false
	}

	if len(in) != 2 && len(in) != 3 {
		return nil, fmt.Errorf("object type must have one or two additional elements, %v", in)
	}

	inMap, ok := in[1].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("object's first additional element must be a map, %v", in[1])
	}

	optionals := make(map[string]bool)
	if len(in) == 3 {
		optionalsSlice, ok := in[2].([]any)
		if !ok {
			return nil, fmt.Errorf("object's second additional element must be a list of strings, %v", in[2])
		}
		for _, val := range optionalsSlice {
			valAsString, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("object's second additional element must be a list of strings, %v", in[2])
			}
			if _, ok := inMap[valAsString]; !ok {
				return nil, fmt.Errorf(
					"an object declared as optional is not in the object itself. This should not be possible:  %v", in,
				)
			}
			optionals[valAsString] = true
		}
	}

	required := []string{}
	properties := make(map[string]any)

	for key, val := range inMap {
		newNode, err := getNodeFromType("", val, false, strict)
		if err != nil {
			return nil, fmt.Errorf("object property %q: %w", key, err)
		}
		properties[key] = newNode
		if !optionals[key] {
			required = append(required, key)
		}
	}

	node["properties"] = properties

	slices.Sort(required)
	node["required"] = required

	return node, nil
}

func getMap(in []any, strict bool) (map[string]any, error) {
	node := map[string]any{
		"type": "object",
	}
	if len(in) != 2 {
		return nil, fmt.Errorf("map type must have exactly one additional element, %v", in)
	}
	newNode, err := getNodeFromType("", in[1], false, strict)
	if err != nil {
		return nil, fmt.Errorf("map: %w", err)
	}
	node["additionalProperties"] = newNode

	return node, nil
}

func getList(in []any, strict bool) (map[string]any, error) {
	node := map[string]any{
		"type": "array",
	}
	if len(in) != 2 {
		return nil, fmt.Errorf("list type must have exactly one additional element, %v", in)
	}

	newNode, err := getNodeFromType("", in[1], false, strict)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	node["items"] = newNode

	return node, nil
}

func getSet(in []any, strict bool) (map[string]any, error) {
	node := map[string]any{
		"type":        "array",
		"uniqueItems": true,
	}
	if len(in) != 2 {
		return nil, fmt.Errorf("set type must have exactly one additional element, %v", in)
	}

	newNode, err := getNodeFromType("", in[1], false, strict)
	if err != nil {
		return nil, fmt.Errorf("set: %w", err)
	}
	node["items"] = newNode

	return node, nil
}

func getTuple(in []any, strict bool) (map[string]any, error) {
	node := map[string]any{
		"type": "array",
	}
	if len(in) != 2 {
		return nil, fmt.Errorf("tuple type must have exactly one additional element, %v", in)
	}

	items := []any{}
	typeSlice, ok := in[1].([]any)
	if !ok {
		return nil, fmt.Errorf("tuple's second argument must be an array, %v", in)
	}

	for _, val := range typeSlice {
		newNode, err := getNodeFromType("", val, false, strict)
		if err != nil {
			return nil, fmt.Errorf("tuple: %w", err)
		}
		items = append(items, newNode)
	}
	node["items"] = items
	node["minItems"] = len(items)
	node["maxItems"] = len(items)

	return node, nil
}
