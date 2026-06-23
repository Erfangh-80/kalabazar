package campaign

import "errors"

var (
	ErrCampaignNotFound       = errors.New("campaign not found")
	ErrCampaignNotApproved    = errors.New("campaign not approved")
	ErrCampaignAlreadyActive  = errors.New("campaign already active")
	ErrCampaignAlreadyApproved = errors.New("campaign already approved")
	ErrCampaignNotStarted     = errors.New("campaign has not started yet")
	ErrCampaignNotExpired     = errors.New("campaign has not expired yet")
)
