package entity_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestNewStore_Success(t *testing.T) {
	phone := "+123456789"
	addr := entity.Address{
		Street: "123 Main St", City: "Tehran", Country: "Iran",
	}
	store, err := entity.NewStore("st-1", "usr-1", "My Store", &phone, &addr, []string{"logo.jpg"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.ID != "st-1" {
		t.Errorf("expected st-1, got %s", store.ID)
	}
	if store.UserID != "usr-1" {
		t.Errorf("expected usr-1, got %s", store.UserID)
	}
	if store.StoreName != "My Store" {
		t.Errorf("expected My Store, got %s", store.StoreName)
	}
	if store.ContactPhone == nil || *store.ContactPhone != "+123456789" {
		t.Errorf("expected +123456789, got %v", store.ContactPhone)
	}
	if store.Address == nil {
		t.Fatal("expected address, got nil")
	}
	if store.Address.Street != "123 Main St" {
		t.Errorf("expected 123 Main St, got %s", store.Address.Street)
	}
	if len(store.MediaAssets) != 1 || store.MediaAssets[0] != "logo.jpg" {
		t.Errorf("expected [logo.jpg], got %v", store.MediaAssets)
	}
}

func TestNewStore_Defaults(t *testing.T) {
	store, err := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.Status != entity.StoreStatusPendingReview {
		t.Errorf("expected pending_review status, got %s", store.Status)
	}
	if !store.IsCommissionApplicable {
		t.Error("expected commission applicable to be true")
	}
	if store.IsBulkSaleEnabled {
		t.Error("expected bulk sale to be false")
	}
}

func TestNewStore_InvalidID(t *testing.T) {
	_, err := entity.NewStore("", "usr-1", "Store", nil, nil, nil)
	if err != entity.ErrStoreInvalidID {
		t.Errorf("expected ErrStoreInvalidID, got %v", err)
	}
}

func TestNewStore_InvalidUserID(t *testing.T) {
	_, err := entity.NewStore("st-1", "", "Store", nil, nil, nil)
	if err != entity.ErrStoreInvalidUserID {
		t.Errorf("expected ErrStoreInvalidUserID, got %v", err)
	}
}

func TestNewStore_InvalidName(t *testing.T) {
	_, err := entity.NewStore("st-1", "usr-1", "", nil, nil, nil)
	if err != entity.ErrStoreInvalidName {
		t.Errorf("expected ErrStoreInvalidName, got %v", err)
	}
}

func TestNewStore_NameTooLong(t *testing.T) {
	name := make([]byte, 256)
	for i := range name {
		name[i] = 'a'
	}
	_, err := entity.NewStore("st-1", "usr-1", string(name), nil, nil, nil)
	if err != entity.ErrStoreNameTooLong {
		t.Errorf("expected ErrStoreNameTooLong, got %v", err)
	}
}

func TestNewStore_InvalidAddress(t *testing.T) {
	addr := entity.Address{Street: "", City: "", Country: ""}
	_, err := entity.NewStore("st-1", "usr-1", "Store", nil, &addr, nil)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestNewStore_EventEmitted(t *testing.T) {
	store, err := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreCreated)
	if !ok {
		t.Fatalf("expected StoreCreated event, got %T", events[0])
	}
	if e.StoreID != "st-1" {
		t.Errorf("expected st-1, got %s", e.StoreID)
	}
}

func TestStore_UpdateInfo_Success(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	phone := "+987654321"
	addr := entity.Address{Street: "456 Oak Ave", City: "Shiraz", Country: "Iran"}

	err := store.UpdateInfo("Updated Store", &phone, &addr, []string{"new.jpg"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.StoreName != "Updated Store" {
		t.Errorf("expected 'Updated Store', got %s", store.StoreName)
	}
	if *store.ContactPhone != "+987654321" {
		t.Errorf("expected +987654321, got %s", *store.ContactPhone)
	}
	if store.Address == nil || store.Address.Street != "456 Oak Ave" {
		t.Errorf("expected 456 Oak Ave, got %v", store.Address)
	}
}

func TestStore_UpdateInfo_InvalidName(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	err := store.UpdateInfo("", nil, nil, nil)
	if err != entity.ErrStoreInvalidName {
		t.Errorf("expected ErrStoreInvalidName, got %v", err)
	}
}

func TestStore_UpdateInfo_InvalidAddress(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	addr := entity.Address{Street: "", City: "", Country: ""}
	err := store.UpdateInfo("Valid Name", nil, &addr, nil)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestStore_Approval_Approve(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	store.Events() // clear initial event

	err := store.Approve()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.Status != entity.StoreStatusActive {
		t.Errorf("expected active, got %s", store.Status)
	}
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.StoreApproved); !ok {
		t.Fatalf("expected StoreApproved, got %T", events[0])
	}
}

func TestStore_Approval_AlreadyApproved(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	store.Events() // clear initial event
	store.Approve()
	store.Events()

	err := store.Approve()
	if err != entity.ErrStoreAlreadyApproved {
		t.Errorf("expected ErrStoreAlreadyApproved, got %v", err)
	}
}

func TestStore_Approval_Reject(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	store.Events()

	err := store.Reject()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.Status != entity.StoreStatusRejected {
		t.Errorf("expected rejected, got %s", store.Status)
	}
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.StoreRejected); !ok {
		t.Fatalf("expected StoreRejected, got %T", events[0])
	}
}

func TestStore_Activate_Success(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	store.Approve()
	store.Deactivate()
	store.Events()

	err := store.Activate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.Status != entity.StoreStatusActive {
		t.Errorf("expected active, got %s", store.Status)
	}
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.StoreActivated); !ok {
		t.Fatalf("expected StoreActivated, got %T", events[0])
	}
}

func TestStore_Activate_FromPendingReview(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	store.Events()

	err := store.Activate()
	if err != entity.ErrStoreNotPendingReview {
		t.Errorf("expected ErrStoreNotPendingReview, got %v", err)
	}
}

func TestStore_Deactivate_Success(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	store.Approve()
	store.Events()

	err := store.Deactivate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.Status != entity.StoreStatusInactive {
		t.Errorf("expected inactive, got %s", store.Status)
	}
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.StoreDeactivated); !ok {
		t.Fatalf("expected StoreDeactivated, got %T", events[0])
	}
}

func TestStore_Deactivate_AlreadyInactive(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	store.Approve()
	store.Deactivate()
	store.Events()

	err := store.Deactivate()
	if err != entity.ErrStoreAlreadyInactive {
		t.Errorf("expected ErrStoreAlreadyInactive, got %v", err)
	}
}

func TestStore_Events_ClearedAfterCall(t *testing.T) {
	store, _ := entity.NewStore("st-1", "usr-1", "My Store", nil, nil, nil)
	events1 := store.Events()
	if len(events1) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events1))
	}
	events2 := store.Events()
	if len(events2) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(events2))
	}
}
