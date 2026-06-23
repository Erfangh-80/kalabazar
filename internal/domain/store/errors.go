package store

import "errors"

var (
	ErrStoreNotFound       = errors.New("store not found")
	ErrStoreAlreadyActive  = errors.New("store is already active")
	ErrStoreSuspended      = errors.New("store is suspended")
)
