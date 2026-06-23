package store

import "time"

type StoreStatus string

const (
	StoreStatusPENDING   StoreStatus = "PENDING"
	StoreStatusACTIVE    StoreStatus = "ACTIVE"
	StoreStatusSUSPENDED StoreStatus = "SUSPENDED"
)

type Store struct {
	ID        int64
	SellerID  int64
	Name      string
	Phone     string
	Status    StoreStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	events    []any
}

func NewStore(sellerID int64, name, phone string) *Store {
	now := time.Now()
	s := &Store{
		SellerID:  sellerID,
		Name:      name,
		Phone:     phone,
		Status:    StoreStatusPENDING,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.emit(StoreCreatedEvent{
		StoreID:  s.ID,
		SellerID: sellerID,
		Name:     name,
	})
	return s
}

func (s *Store) Activate() error {
	if s.Status == StoreStatusACTIVE {
		return ErrStoreAlreadyActive
	}
	if s.Status == StoreStatusSUSPENDED {
		return ErrStoreSuspended
	}
	s.Status = StoreStatusACTIVE
	s.UpdatedAt = time.Now()
	s.emit(StoreActivatedEvent{
		StoreID: s.ID,
	})
	return nil
}

func (s *Store) Suspend() {
	s.Status = StoreStatusSUSPENDED
	s.UpdatedAt = time.Now()
}

func (s *Store) Events() []any {
	evts := s.events
	s.events = nil
	return evts
}

func (s *Store) emit(event any) {
	s.events = append(s.events, event)
}
