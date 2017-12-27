package types

// Case defines a test case
type Case struct {
	// Description describe
	Description string `json:"description,omitempty"`

	// Flow defines test flow of a test case
	Flow []RoundTrip `json:"flow,omitempty"`
}
