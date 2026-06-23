package commission_test

import (
	"testing"

	"stock-service-version-three/internal/application/commission"
	domaincommission "stock-service-version-three/internal/domain/commission"
)

type mockCommissionRepo struct {
	commissions map[int64]*domaincommission.Commission
}

func newMockCommissionRepo() *mockCommissionRepo {
	return &mockCommissionRepo{commissions: make(map[int64]*domaincommission.Commission)}
}

func (m *mockCommissionRepo) Save(c *domaincommission.Commission) error {
	m.commissions[c.ID()] = c
	return nil
}

func (m *mockCommissionRepo) FindByID(id int64) (*domaincommission.Commission, error) {
	c, ok := m.commissions[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *mockCommissionRepo) FindBySellerID(sellerID int64) ([]*domaincommission.Commission, error) {
	var result []*domaincommission.Commission
	for _, c := range m.commissions {
		if c.SellerID() == sellerID {
			result = append(result, c)
		}
	}
	return result, nil
}

func TestCalculateCommission_Success(t *testing.T) {
	repo := newMockCommissionRepo()
	uc := commission.NewCalculateCommissionUseCase(repo)

	req := commission.CalculateCommissionRequest{
		SellerID:    1,
		SalesAmount: 2040000,
		Rate:        0.10,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.CommissionID == 0 {
		t.Error("expected non-zero CommissionID")
	}
	if resp.SalesAmount != req.SalesAmount {
		t.Errorf("expected SalesAmount %d, got %d", req.SalesAmount, resp.SalesAmount)
	}
	expectedAmount := int64(float64(req.SalesAmount) * req.Rate)
	if resp.Amount != expectedAmount {
		t.Errorf("expected Amount %d, got %d", expectedAmount, resp.Amount)
	}

	if len(repo.commissions) != 1 {
		t.Errorf("expected 1 commission saved, got %d", len(repo.commissions))
	}
}

func TestCalculateCommission_InvalidRateZero(t *testing.T) {
	repo := newMockCommissionRepo()
	uc := commission.NewCalculateCommissionUseCase(repo)

	req := commission.CalculateCommissionRequest{
		SellerID:    1,
		SalesAmount: 2040000,
		Rate:        0,
	}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatal("expected error for zero rate, got nil")
	}
	if err != domaincommission.ErrInvalidCommissionRate {
		t.Errorf("expected ErrInvalidCommissionRate, got %v", err)
	}
}

func TestCalculateCommission_InvalidRateNegative(t *testing.T) {
	repo := newMockCommissionRepo()
	uc := commission.NewCalculateCommissionUseCase(repo)

	req := commission.CalculateCommissionRequest{
		SellerID:    1,
		SalesAmount: 2040000,
		Rate:        -0.05,
	}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
	if err != domaincommission.ErrInvalidCommissionRate {
		t.Errorf("expected ErrInvalidCommissionRate, got %v", err)
	}
}

func TestCalculateCommission_ZeroSalesAmount(t *testing.T) {
	repo := newMockCommissionRepo()
	uc := commission.NewCalculateCommissionUseCase(repo)

	req := commission.CalculateCommissionRequest{
		SellerID:    1,
		SalesAmount: 0,
		Rate:        0.10,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Amount != 0 {
		t.Errorf("expected Amount 0, got %d", resp.Amount)
	}
	if len(repo.commissions) != 1 {
		t.Errorf("expected 1 commission saved, got %d", len(repo.commissions))
	}
}
