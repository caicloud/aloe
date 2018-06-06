package types

// PresetConfig defines config of presetter
type PresetConfig struct {
	// Name defines name of presetter
	Name string `json:"name"`

	// Args defines preset args
	Args map[string]Template `json:"args,omitempty"`
}
