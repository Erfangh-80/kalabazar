package payout_test

import (
	"testing"

	"stock-service-version-three/internal/application/payout"
	domainpayout "stock-service-version-three/internal/domain/payout"
)

type mockPayoutRepo struct {
	payouts map[int64]*domainpayout.Payout
}

func newMockPayoutRepo() *mockPayoutRepo {
	return &mockPayoutRepo{payouts: make(map[int64]*domainpayout.Payout)}
}

func (m *mockPayoutRepo) Save(p *domainpayout.Payout) error {
	m.payouts[p.ID()] = p
	return nil
}

func (m *mockPayoutRepo) FindByID(id int64) (*domainpayout.Payout, error) {
	p, ok := m.payouts[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockPayoutRepo) FindBySellerID(sellerID int64) ([]*domainpayout.Payout, error) {
	var result []*domainpayout.Payout
	for _, p := range m.payouts {
		if p.SellerID() == sellerID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockPayoutRepo) Update(p *domainpayout.Payout) error {
	m.payouts[p.ID()] = p
	return nil
}

func TestExecutePayout_Success(t *testing.T) {
	repo := newMockPayoutRepo()
	uc := payout.NewExecutePayoutUseCase(repo)

	req := payout.ExecutePayoutRequest{
		SellerID: 1,
		Amount:   1836000,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.PayoutID == 0 {
		t.Error("expected non-zero PayoutID")
	}
	if resp.Amount != req.Amount {
		t.Errorf("expected Amount %d, got %d", req.Amount, resp.Amount)
	}
	if resp.Status != domainpayout.PayoutStatusExecuted {
		t.Errorf("expected status %s, got %s", domainpayout.PayoutStatusExecuted, resp.Status)
	}
	if len(repo.payouts) != 1 {
		t.Errorf("expected 1 payout saved, got %d", len(repo.payouts))
	}
}

func TestExecutePayout_ZeroAmount(t *testing.T) {
	repo := newMockPayoutRepo()
	uc := payout.NewExecutePayoutUseCase(repo)

	req := payout.ExecutePayoutRequest{
		SellerID: 1,
		Amount:   0,
	}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatal("expected error for zero amount, got nil")
	}
	if err != domainpayout.ErrInvalidPayoutAmount {
		t.Errorf("expected ErrInvalidPayoutAmount, got %v", err)
	}
}

func TestExecutePayout_NegativeAmount(t *testing.T) {
	repo := newMockPayoutRepo()
	uc := payout.NewExecutePayoutUseCase(repo)

	req := payout.ExecutePayoutRequest{
		SellerID: 1,
		Amount:   -500,
	}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatal("expected error for negative amount, got nil")
	}
	if err != domainpayout.ErrInvalidPayoutAmount {
		t.Errorf("expected ErrInvalidPayoutAmount, got %v", err)
	}
}
