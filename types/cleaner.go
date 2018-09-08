package types

// CleanerConfig defines config of cleaner
type CleanerConfig struct {
	// Name defines name of cleaner
	Name string `json:"name"`

	// ForEach defines whether cleaner will clean context for each cases
	ForEach bool `json:"forEach"`

	// Args defines cleaner args
	Args map[string]Template `json:"args,omitempty"`
}
