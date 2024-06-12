package jsonschema

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateSchema(t *testing.T) {
	t.Parallel()
	tfPath := "../../test/modules"
	schemaPath := "../../test/expected"
	testCases := []string{
		"empty",
		"simple",
		"simple-types",
		"complex-types",
		"custom-validation",
	}
	for i := range testCases {
		name := testCases[i]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			expected, err := os.ReadFile(filepath.Join(schemaPath, name, "schema.json"))
			require.NoError(t, err)

			result, err := CreateSchema(filepath.Join(tfPath, name), false)
			require.NoError(t, err)

			var expectedMap, resultMap map[string]interface{}
			err = json.Unmarshal(expected, &expectedMap)
			require.NoError(t, err)

			err = json.Unmarshal([]byte(result), &resultMap)
			require.NoError(t, err)

			if d := cmp.Diff(expectedMap, resultMap); d != "" {
				t.Errorf("Schema has incorrect value (-want,+got):\n%s", d)
			}
		})
	}
}

func TestSampleInput(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		schemaPath       string
		filePath         string
		keywordLocations []string
	}{
		// {
		// 	name:       "empty",
		// 	filePath:   "test/expected/empty/sample-input/test-input-min.json",
		// 	schemaPath: "test/expected/empty/schema.json",
		// 	err:        nil,
		// },
		// {
		// 	name:       "simple",
		// 	filePath:   "test/expected/simple/sample-input/test-input-all.json",
		// 	schemaPath: "test/expected/simple/schema.json",
		// 	causes:     nil,
		// },
		// {
		// 	name:       "simple-types",
		// 	filePath:   "test/expected/simple-types/sample-input/test-input-min.json",
		// 	schemaPath: "test/expected/simple-types/schema.json",
		// 	err:        nil,
		// },
		// {
		// 	name:       "complex-types",
		// 	filePath:   "test/expected/complex-types/sample-input/test-input-min.json",
		// 	schemaPath: "test/expected/complex-types/schema.json",
		// 	err:        nil,
		// },
		{
			name:       "custom-validation",
			filePath:   "../../test/expected/custom-validation/sample-input/test-input-bad.json",
			schemaPath: "../../test/expected/custom-validation/schema.json",
			keywordLocations: []string{
				"/properties/a_list_maximum_minimum_length/minItems",
				"/properties/a_map_maximum_minimum_entries/minProperties",
				"/properties/a_number_enum_kind_1/type",
				"/properties/a_number_enum_kind_2/enum",
				"/properties/a_number_exclusive_maximum_minimum/exclusiveMaximum",
				"/properties/a_number_maximum_minimum/maximum",
				"/properties/a_set_maximum_minimum_items",
				"/properties/a_string_enum_escaped_characters_kind_1/enum",
				"/properties/a_string_enum_escaped_characters_kind_2/enum",
				"/properties/a_string_enum_kind_1/enum",
				"/properties/a_string_enum_kind_2/type",
				"/properties/a_string_maximum_minimum_length/maxLength",
				"/properties/a_string_pattern_1/pattern",
				"/properties/a_string_pattern_2/pattern",
				"/properties/an_object_maximum_minimum_items",
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path, err := filepath.Abs(tc.schemaPath)
			require.NoError(t, err)
			c := jsonschema.NewCompiler()
			f, err := os.Open(path)
			require.NoError(t, err)
			err = c.AddResource("file://"+path, f)
			require.NoError(t, err)
			s, err := c.Compile("file://" + path)
			require.NoError(t, err)

			input, err := os.ReadFile(tc.filePath)
			require.NoError(t, err)
			var m interface{}
			err = json.Unmarshal(input, &m)
			require.NoError(t, err)

			var keywordLocations []string
			err = s.Validate(m)
			if err != nil {
				valErr := &jsonschema.ValidationError{}
				ok := errors.As(err, &valErr)
				require.True(t, ok)
				for _, cause := range valErr.Causes {
					keywordLocations = append(keywordLocations, cause.KeywordLocation)
				}
			}
			slices.Sort(keywordLocations)
			require.Equal(t, tc.keywordLocations, keywordLocations)
		})
	}
}
