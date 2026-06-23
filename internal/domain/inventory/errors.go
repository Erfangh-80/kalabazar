package inventory

import "errors"

var (
	ErrInventoryNotFound       = errors.New("inventory not found")
	ErrInsufficientStock       = errors.New("insufficient available stock")
	ErrInsufficientReservedStock = errors.New("insufficient reserved stock")
	ErrInvalidPrice            = errors.New("price must be greater than zero")
	ErrInvalidStock            = errors.New("stock must not be negative")
	ErrInvalidQuantity         = errors.New("quantity must not be negative")
)
