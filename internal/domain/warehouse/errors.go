package warehouse

import "errors"

var (
	ErrWarehouseNotFound = errors.New("warehouse not found")
	ErrInvalidCapacity   = errors.New("invalid capacity: must be greater than zero")
)
