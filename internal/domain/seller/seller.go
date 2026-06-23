package seller

import "time"

type SellerStatus string

const (
	SellerStatusUnverified SellerStatus = "UNVERIFIED"
	SellerStatusVerified   SellerStatus = "VERIFIED"
)

type Seller struct {
	ID        int64
	UserID    int64
	Name      string
	Phone     string
	Status    SellerStatus
	Score     float64
	Rank      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSeller(userID int64, name string, phone string) *Seller {
	now := time.Now()
	return &Seller{
		UserID:    userID,
		Name:      name,
		Phone:     phone,
		Status:    SellerStatusUnverified,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *Seller) VerifyKYC() error {
	if s.Status == SellerStatusVerified {
		return ErrSellerAlreadyVerified
	}
	s.Status = SellerStatusVerified
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Seller) UpdateRank(score float64, rank string) {
	s.Score = score
	s.Rank = rank
	s.UpdatedAt = time.Now()
}
