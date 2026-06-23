package settlement

import "errors"

var ErrInvalidSettlementAmount = errors.New("commission cannot exceed gross sales and amounts must be non-negative")
