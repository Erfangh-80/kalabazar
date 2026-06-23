package settlement_test

import (
	"testing"

	"stock-service-version-three/internal/application/settlement"
	domainsettlement "stock-service-version-three/internal/domain/settlement"
)

type mockSettlementRepo struct {
	settlements map[int64]*domainsettlement.Settlement
}

func newMockSettlementRepo() *mockSettlementRepo {
	return &mockSettlementRepo{settlements: make(map[int64]*domainsettlement.Settlement)}
}

func (m *mockSettlementRepo) Save(s *domainsettlement.Settlement) error {
	m.settlements[s.ID()] = s
	return nil
}

func (m *mockSettlementRepo) FindByID(id int64) (*domainsettlement.Settlement, error) {
	s, ok := m.settlements[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (m *mockSettlementRepo) FindBySellerID(sellerID int64) ([]*domainsettlement.Settlement, error) {
	var result []*domainsettlement.Settlement
	for _, s := range m.settlements {
		if s.SellerID() == sellerID {
			result = append(result, s)
		}
	}
	return result, nil
}

func TestCreateSettlement_Success(t *testing.T) {
	repo := newMockSettlementRepo()
	uc := settlement.NewCreateSettlementUseCase(repo)

	req := settlement.CreateSettlementRequest{
		SellerID:   1,
		GrossSales: 2040000,
		Commission: 204000,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.SettlementID == 0 {
		t.Error("expected non-zero SettlementID")
	}
	if resp.GrossSales != req.GrossSales {
		t.Errorf("expected GrossSales %d, got %d", req.GrossSales, resp.GrossSales)
	}
	if resp.Commission != req.Commission {
		t.Errorf("expected Commission %d, got %d", req.Commission, resp.Commission)
	}
	expectedNet := req.GrossSales - req.Commission
	if resp.NetAmount != expectedNet {
		t.Errorf("expected NetAmount %d, got %d", expectedNet, resp.NetAmount)
	}

	if len(repo.settlements) != 1 {
		t.Errorf("expected 1 settlement saved, got %d", len(repo.settlements))
	}
}

func TestCreateSettlement_CommissionExceedsGrossSales(t *testing.T) {
	repo := newMockSettlementRepo()
	uc := settlement.NewCreateSettlementUseCase(repo)

	req := settlement.CreateSettlementRequest{
		SellerID:   1,
		GrossSales: 1000000,
		Commission: 2000000,
	}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatal("expected error when commission exceeds gross sales, got nil")
	}
	if err != domainsettlement.ErrInvalidSettlementAmount {
		t.Errorf("expected ErrInvalidSettlementAmount, got %v", err)
	}
}

func TestCreateSettlement_NegativeGrossSales(t *testing.T) {
	repo := newMockSettlementRepo()
	uc := settlement.NewCreateSettlementUseCase(repo)

	req := settlement.CreateSettlementRequest{
		SellerID:   1,
		GrossSales: -100,
		Commission: 10,
	}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatal("expected error for negative gross sales, got nil")
	}
	if err != domainsettlement.ErrInvalidSettlementAmount {
		t.Errorf("expected ErrInvalidSettlementAmount, got %v", err)
	}
}

func TestCreateSettlement_NegativeCommission(t *testing.T) {
	repo := newMockSettlementRepo()
	uc := settlement.NewCreateSettlementUseCase(repo)

	req := settlement.CreateSettlementRequest{
		SellerID:   1,
		GrossSales: 1000000,
		Commission: -50,
	}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatal("expected error for negative commission, got nil")
	}
	if err != domainsettlement.ErrInvalidSettlementAmount {
		t.Errorf("expected ErrInvalidSettlementAmount, got %v", err)
	}
}

func TestCreateSettlement_ZeroValues(t *testing.T) {
	repo := newMockSettlementRepo()
	uc := settlement.NewCreateSettlementUseCase(repo)

	req := settlement.CreateSettlementRequest{
		SellerID:   1,
		GrossSales: 0,
		Commission: 0,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("expected no error for zero values, got %v", err)
	}
	if resp.NetAmount != 0 {
		t.Errorf("expected NetAmount 0, got %d", resp.NetAmount)
	}
}
