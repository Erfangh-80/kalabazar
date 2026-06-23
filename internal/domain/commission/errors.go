package commission

import "errors"

var ErrInvalidCommissionRate = errors.New("commission rate must be greater than zero")
