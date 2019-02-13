package types

// Var defines variable
type Var struct {
	// Name defines variable name
	Name string `json:"name"`

	// Selector select variable value from response
	Selector []Template `json:"selector"`
}

// Definition defines new variable from response
type Definition struct {
	Var `json:",inline"`
	// Type defines variable from
	// enum ["body", "status", "header"]
	// default is body
	Type string `json:"type,omitempty"`
}
