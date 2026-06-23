package seller

import "errors"

var (
	ErrSellerNotFound         = errors.New("seller not found")
	ErrInvalidKYCStatus       = errors.New("invalid KYC status")
	ErrSellerAlreadyVerified  = errors.New("seller already verified")
	ErrInvalidSellerName      = errors.New("invalid seller name")
	ErrInvalidPhone           = errors.New("invalid phone number")
)
