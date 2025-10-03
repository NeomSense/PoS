package types

// Params defines the parameters for the blog module
type Params struct {
	// Add any module parameters here, e.g., MaxPostLength uint64
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{}
}

// Validate validates the parameters
func (p Params) Validate() error {
	return nil
}
