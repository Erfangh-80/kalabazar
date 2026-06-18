package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"kalabazar-stock-service/internal/domain/event"
)

var (
	ErrStoreInvalidName     = errors.New("store name cannot be empty")
	ErrStoreNameTooLong     = errors.New("store name cannot exceed 255 characters")
	ErrStoreAlreadyActive   = errors.New("store is already active")
	ErrStoreAlreadyInactive = errors.New("store is already inactive")
	ErrStoreInvalidID       = errors.New("store id cannot be empty")
	ErrStoreInvalidUserID   = errors.New("user id cannot be empty")
	ErrStoreNotFound        = errors.New("store not found")
)

// StoreStatus represents the operational status of a store.
type StoreStatus string

const (
	StoreStatusActive   StoreStatus = "active"
	StoreStatusInactive StoreStatus = "inactive"
)

// Store represents a seller's store in the marketplace.
type Store struct {
	ID                     string
	UserID                 string
	StoreName              string
	ContactPhone           *string
	Address                *Address
	MediaAssets            []string
	Status                 StoreStatus
	IsCommissionApplicable bool
	IsBulkSaleEnabled      bool
	CreatedAt              time.Time
	UpdatedAt              time.Time

	events []any
}

// NewStore creates a new Store with active status and default values.
func NewStore(id, userID, storeName string, contactPhone *string, address *Address, mediaAssets []string) (*Store, error) {
	now := time.Now()
	store := &Store{
		ID:                     id,
		UserID:                 userID,
		StoreName:              storeName,
		ContactPhone:           contactPhone,
		Address:                address,
		MediaAssets:            mediaAssets,
		Status:                 StoreStatusActive,
		IsCommissionApplicable: true,
		IsBulkSaleEnabled:      false,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if err := store.validate(); err != nil {
		return nil, err
	}
	store.events = append(store.events, event.StoreCreated{
		StoreID:   id,
		UserID:    userID,
		StoreName: storeName,
		Timestamp: now,
	})
	return store, nil
}

// validate checks all invariant business rules for the Store entity.
func (s *Store) validate() error {
	switch {
	case s.ID == "":
		return ErrStoreInvalidID
	case s.UserID == "":
		return ErrStoreInvalidUserID
	case s.StoreName == "":
		return ErrStoreInvalidName
	case utf8.RuneCountInString(s.StoreName) > 255:
		return ErrStoreNameTooLong
	case s.Address != nil && s.Address.Validate() != nil:
		return s.Address.Validate()
	default:
		return nil
	}
}

// UpdateInfo updates the store's mutable information fields.
func (s *Store) UpdateInfo(storeName string, contactPhone *string, address *Address, mediaAssets []string) error {
	oldName, oldAddr := s.StoreName, s.Address
	s.StoreName = storeName
	s.Address = address
	if err := s.validate(); err != nil {
		s.StoreName, s.Address = oldName, oldAddr
		return err
	}

	s.ContactPhone = contactPhone
	s.MediaAssets = mediaAssets
	s.UpdatedAt = time.Now()
	s.events = append(s.events, event.StoreUpdated{
		StoreID:   s.ID,
		Timestamp: s.UpdatedAt,
	})
	return nil
}

// Activate transitions the store to the active status.
func (s *Store) Activate() error {
	if s.Status == StoreStatusActive {
		return ErrStoreAlreadyActive
	}
	s.Status = StoreStatusActive
	s.UpdatedAt = time.Now()
	s.events = append(s.events, event.StoreActivated{
		StoreID:   s.ID,
		Timestamp: s.UpdatedAt,
	})
	return nil
}

// Deactivate transitions the store to the inactive status.
func (s *Store) Deactivate() error {
	if s.Status == StoreStatusInactive {
		return ErrStoreAlreadyInactive
	}
	s.Status = StoreStatusInactive
	s.UpdatedAt = time.Now()
	s.events = append(s.events, event.StoreDeactivated{
		StoreID:   s.ID,
		Timestamp: s.UpdatedAt,
	})
	return nil
}

// Events returns and clears the domain events produced by the entity.
func (s *Store) Events() []any {
	events := s.events
	s.events = nil
	return events
}

type StoreRepository interface {
	Save(store *Store) error
	FindByID(id string) (*Store, error)
	FindByUserID(userID string) ([]*Store, error)
	Update(store *Store) error
}
