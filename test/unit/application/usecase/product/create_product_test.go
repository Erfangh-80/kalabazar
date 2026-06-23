package product_test

import (
	"errors"
	"testing"

	appProduct "stock-service-version-three/internal/application/product"
	domain "stock-service-version-three/internal/domain/product"
)

type mockProductRepository struct {
	products map[int64]*domain.Product
	nextID   int64
	saveErr  error
	findErr  error
}

func (m *mockProductRepository) Save(p *domain.Product) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.nextID++
	p.ID = m.nextID
	m.products[p.ID] = p
	return nil
}

func (m *mockProductRepository) FindByID(id int64) (*domain.Product, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	p, ok := m.products[id]
	if !ok {
		return nil, domain.ErrProductNotFound
	}
	return p, nil
}

func (m *mockProductRepository) FindByStoreID(storeID int64) ([]*domain.Product, error) {
	return nil, nil
}

func (m *mockProductRepository) Update(p *domain.Product) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.products[p.ID] = p
	return nil
}

func TestCreateProduct_Success(t *testing.T) {
	repo := &mockProductRepository{products: make(map[int64]*domain.Product)}
	uc := appProduct.NewCreateProductUseCase(repo)

	resp, err := uc.Execute(appProduct.CreateProductRequest{
		StoreID:    1,
		Title:      "Runner 3000",
		CategoryID: 12,
		Brand:      "SportLine",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ProductID != 1 {
		t.Errorf("expected ProductID 1, got %d", resp.ProductID)
	}
	if resp.Status != string(domain.PENDING_REVIEW) {
		t.Errorf("expected Status %s, got %s", domain.PENDING_REVIEW, resp.Status)
	}
}

func TestCreateProduct_EmptyTitle(t *testing.T) {
	repo := &mockProductRepository{products: make(map[int64]*domain.Product)}
	uc := appProduct.NewCreateProductUseCase(repo)

	_, err := uc.Execute(appProduct.CreateProductRequest{
		StoreID:    1,
		Title:      "",
		CategoryID: 12,
		Brand:      "SportLine",
	})
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestCreateProduct_RepositoryError(t *testing.T) {
	expectedErr := errors.New("db error")
	repo := &mockProductRepository{
		products: make(map[int64]*domain.Product),
		saveErr:  expectedErr,
	}
	uc := appProduct.NewCreateProductUseCase(repo)

	_, err := uc.Execute(appProduct.CreateProductRequest{
		StoreID:    1,
		Title:      "Runner 3000",
		CategoryID: 12,
		Brand:      "SportLine",
	})
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}
