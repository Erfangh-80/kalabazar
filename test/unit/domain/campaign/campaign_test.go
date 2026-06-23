package campaign_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/domain/campaign"
)

func TestNewCampaign(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Launch Discount", "percentage", 15, now, now.Add(24*time.Hour))

	if c.Title != "Launch Discount" {
		t.Errorf("expected title 'Launch Discount', got %s", c.Title)
	}
	if c.DiscountType != "percentage" {
		t.Errorf("expected discount_type 'percentage', got %s", c.DiscountType)
	}
	if c.Value != 15 {
		t.Errorf("expected value 15, got %f", c.Value)
	}
	if c.Status != campaign.CampaignStatusInactive {
		t.Errorf("expected status INACTIVE, got %s", c.Status)
	}
	if c.ApprovalStatus != campaign.ApprovalStatusPending {
		t.Errorf("expected approval_status PENDING, got %s", c.ApprovalStatus)
	}
	if c.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestCampaignApprove(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))

	events, err := c.Approve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ApprovalStatus != campaign.ApprovalStatusApproved {
		t.Errorf("expected APPROVED, got %s", c.ApprovalStatus)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	evt, ok := events[0].(campaign.CampaignApprovedEvent)
	if !ok {
		t.Fatalf("expected CampaignApprovedEvent, got %T", events[0])
	}
	if evt.CampaignID != c.ID {
		t.Errorf("expected campaign ID %d, got %d", c.ID, evt.CampaignID)
	}
}

func TestCampaignApprove_AlreadyApproved(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	c.Approve()

	_, err := c.Approve()
	if err != campaign.ErrCampaignAlreadyApproved {
		t.Errorf("expected ErrCampaignAlreadyApproved, got %v", err)
	}
}

func TestCampaignActivate_BeforeStart(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now.Add(1*time.Hour), now.Add(24*time.Hour))
	c.Approve()

	_, err := c.Activate(now)
	if err != campaign.ErrCampaignNotStarted {
		t.Errorf("expected ErrCampaignNotStarted, got %v", err)
	}
	if c.IsActive() {
		t.Error("expected campaign not to be active")
	}
}

func TestCampaignActivate_NotApproved(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))

	_, err := c.Activate(now)
	if err != campaign.ErrCampaignNotApproved {
		t.Errorf("expected ErrCampaignNotApproved, got %v", err)
	}
}

func TestCampaignActivate_Success(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	c.Approve()

	events, err := c.Activate(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Status != campaign.CampaignStatusActive {
		t.Errorf("expected ACTIVE, got %s", c.Status)
	}
	if !c.IsActive() {
		t.Error("expected IsActive to be true")
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	evt, ok := events[0].(campaign.CampaignActivatedEvent)
	if !ok {
		t.Fatalf("expected CampaignActivatedEvent, got %T", events[0])
	}
	if evt.CampaignID != c.ID {
		t.Errorf("expected campaign ID %d, got %d", c.ID, evt.CampaignID)
	}
}

func TestCampaignActivate_AlreadyActive(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	c.Approve()
	c.Activate(now)

	_, err := c.Activate(now)
	if err != campaign.ErrCampaignAlreadyActive {
		t.Errorf("expected ErrCampaignAlreadyActive, got %v", err)
	}
}

func TestCampaignEnd(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	c.Approve()
	c.Activate(now)

	events, err := c.End(now.Add(48 * time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Status != campaign.CampaignStatusInactive {
		t.Errorf("expected INACTIVE, got %s", c.Status)
	}
	if c.IsActive() {
		t.Error("expected IsActive to be false")
	}
	if !c.IsExpiredAt(now.Add(48 * time.Hour)) {
		t.Error("expected IsExpired to be true")
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	evt, ok := events[0].(campaign.CampaignEndedEvent)
	if !ok {
		t.Fatalf("expected CampaignEndedEvent, got %T", events[0])
	}
	if evt.CampaignID != c.ID {
		t.Errorf("expected campaign ID %d, got %d", c.ID, evt.CampaignID)
	}
}

func TestCampaignEnd_BeforeEnd(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))

	_, err := c.End(now)
	if err != campaign.ErrCampaignNotExpired {
		t.Errorf("expected ErrCampaignNotExpired, got %v", err)
	}
}

func TestIsActive(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))

	if c.IsActive() {
		t.Error("expected new campaign to not be active")
	}

	c.Approve()
	c.Activate(now)
	if !c.IsActive() {
		t.Error("expected activated campaign to be active")
	}
}

func TestIsExpired(t *testing.T) {
	now := time.Now()
	endAt := now.Add(24 * time.Hour)
	c := campaign.NewCampaign("Test", "percentage", 10, now, endAt)

	if c.IsExpired() {
		t.Error("expected new campaign to not be expired")
	}

	if !c.IsExpiredAt(now.Add(48 * time.Hour)) {
		t.Error("expected campaign to be expired after end_at")
	}
}

func TestNewCampaignInventoryLink(t *testing.T) {
	link := campaign.NewCampaignInventoryLink(1, 2)

	if link.CampaignID != 1 {
		t.Errorf("expected CampaignID 1, got %d", link.CampaignID)
	}
	if link.InventoryID != 2 {
		t.Errorf("expected InventoryID 2, got %d", link.InventoryID)
	}
	if link.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestLinkToInventory(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))

	events := c.LinkToInventory(42)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	evt, ok := events[0].(campaign.CampaignLinkedToInventoryEvent)
	if !ok {
		t.Fatalf("expected CampaignLinkedToInventoryEvent, got %T", events[0])
	}
	if evt.CampaignID != c.ID {
		t.Errorf("expected CampaignID %d, got %d", c.ID, evt.CampaignID)
	}
	if evt.InventoryID != 42 {
		t.Errorf("expected InventoryID 42, got %d", evt.InventoryID)
	}
}

func TestNewCampaign_EmitsCreatedEvent(t *testing.T) {
	now := time.Now()
	c := campaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))

	events := c.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	evt, ok := events[0].(campaign.CampaignCreatedEvent)
	if !ok {
		t.Fatalf("expected CampaignCreatedEvent, got %T", events[0])
	}
	if evt.Title != "Test" {
		t.Errorf("expected title 'Test', got %s", evt.Title)
	}
}

func TestIDAssignment(t *testing.T) {
	now := time.Now()
	c1 := campaign.NewCampaign("A", "percentage", 10, now, now.Add(24*time.Hour))
	c2 := campaign.NewCampaign("B", "percentage", 20, now, now.Add(24*time.Hour))

	if c1.ID == c2.ID {
		t.Error("expected different IDs for different campaigns")
	}
	if c1.ID <= 0 {
		t.Error("expected positive ID")
	}
}
