package inventory

import "errors"

var (
	ErrInvalidReferencePrice  = errors.New("reference price must be greater than zero")
	ErrInvalidReferenceSource = errors.New("reference source must not be empty")
)
