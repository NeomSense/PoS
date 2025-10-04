package types

// DefaultParams returns default module parameters.
func DefaultParams() Params {
	return Params{}
}

// Validate performs basic validation of module parameters.
func (p Params) Validate() error {
	return nil
}
