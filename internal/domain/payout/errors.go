package payout

import "errors"

var (
	ErrPayoutNotFound        = errors.New("payout not found")
	ErrPayoutAlreadyExecuted = errors.New("payout has already been executed")
	ErrInvalidPayoutAmount   = errors.New("payout amount must be greater than zero")
)
