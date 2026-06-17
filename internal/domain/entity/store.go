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
	if id == "" {
		return nil, ErrStoreInvalidID
	}
	if userID == "" {
		return nil, ErrStoreInvalidUserID
	}
	if storeName == "" {
		return nil, ErrStoreInvalidName
	}
	if utf8.RuneCountInString(storeName) > 255 {
		return nil, ErrStoreNameTooLong
	}
	if address != nil {
		if err := address.Validate(); err != nil {
			return nil, err
		}
	}

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
	store.events = append(store.events, event.StoreCreated{
		StoreID:   id,
		UserID:    userID,
		StoreName: storeName,
		Timestamp: now,
	})
	return store, nil
}

// UpdateInfo updates the store's mutable information fields.
func (s *Store) UpdateInfo(storeName string, contactPhone *string, address *Address, mediaAssets []string) error {
	if storeName == "" {
		return ErrStoreInvalidName
	}
	if utf8.RuneCountInString(storeName) > 255 {
		return ErrStoreNameTooLong
	}
	if address != nil {
		if err := address.Validate(); err != nil {
			return err
		}
	}

	s.StoreName = storeName
	s.ContactPhone = contactPhone
	s.Address = address
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
