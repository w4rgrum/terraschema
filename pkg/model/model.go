package model

import (
	"github.com/hashicorp/hcl/v2"
)

// Each variable block in the Terraform configuration file is marshalled into this struct.
type Variable struct {
	Default     hcl.Expression `hcl:"default"`
	Description *string        `hcl:"description"`
	Nullable    *bool          `hcl:"nullable"`
	// Sensitive is ignored.
	Sensitive *bool `hcl:"sensitive"`
	// Validation blocks can be used to add extra rules to the JSON schema, as long as their conditions are written in a certain format.
	Validation *ValidationBlock `hcl:"validation,block"`
	Type       hcl.Expression   `hcl:"type"`
}

type ValidationBlock struct {
	Condition hcl.Expression `hcl:"condition"`
	// ErrorMessage is ignored.
	ErrorMessage string `hcl:"error_message"`
}

// TranslatedVariable contains the Variable struct, as well as some extra information that can be used for debugging.
// This is done here because it can be difficult to extract this information from pure hcl.Expressions as present in
// the Variable struct without further context. Required is used internally, and the others are for debugging.
type TranslatedVariable struct {
	// if the variable has a validation block, this stores its condition as a string. This is useful for debugging.
	ConditionAsString *string
	// DefaultAsString is the default value of the variable, as a string. This is useful for debugging complex default values.
	DefaultAsString *string
	// TypeAsString is the type of the variable, as a string. This is useful for debugging complex types, such as objects.
	// A nil value for type implies no type has been specified, so the terraform code accepts "any" type.
	TypeAsString *string
	// Required is true if and only if the variable has no default value.
	Required bool
	// The variable block used to generate the other fields in this struct.
	Variable Variable
}
