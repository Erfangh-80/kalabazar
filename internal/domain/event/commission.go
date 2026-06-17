package event

import "time"

// CommissionRuleCreated is emitted when a new commission rule is defined.
type CommissionRuleCreated struct {
	CommissionID string
	ProductID    string
	RatePercent  float64
	Timestamp    time.Time
}

func (e CommissionRuleCreated) EventName() string { return "commission.rule.created" }

// CommissionCalculated is emitted when a commission amount is calculated for a sale.
type CommissionCalculated struct {
	CommissionID string
	SaleAmount   float64
	CommissionAmount float64
	Timestamp    time.Time
}

func (e CommissionCalculated) EventName() string { return "commission.calculated" }
