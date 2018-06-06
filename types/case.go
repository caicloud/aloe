package types

// Case defines a test case
type Case struct {
	// Summary describes the test case
	Summary string `json:"summary,omitempty"`

	// Flow defines test flow of a test case
	Flow []RoundTrip `json:"flow,omitempty"`
}
