package product

import "errors"

var (
	ErrProductNotFound        = errors.New("product not found")
	ErrProductAlreadyApproved = errors.New("product already approved")
)
