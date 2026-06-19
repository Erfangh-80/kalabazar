package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockInventoryStoreRepo struct {
	stores map[string]*entity.Store
}

func (m *mockInventoryStoreRepo) Save(store *entity.Store) error {
	m.stores[store.ID] = store
	return nil
}

func (m *mockInventoryStoreRepo) FindByID(id string) (*entity.Store, error) {
	s, ok := m.stores[id]
	if !ok {
		return nil, entity.ErrStoreNotFound
	}
	return s, nil
}

func (m *mockInventoryStoreRepo) FindByUserID(userID string) ([]*entity.Store, error) {
	return nil, nil
}

func (m *mockInventoryStoreRepo) Update(store *entity.Store) error {
	m.stores[store.ID] = store
	return nil
}

type mockInventoryWarehouseRepo struct {
	warehouses map[string]*entity.Warehouse
}

func (m *mockInventoryWarehouseRepo) Save(warehouse *entity.Warehouse) error {
	m.warehouses[warehouse.ID] = warehouse
	return nil
}

func (m *mockInventoryWarehouseRepo) FindByID(id string) (*entity.Warehouse, error) {
	w, ok := m.warehouses[id]
	if !ok {
		return nil, entity.ErrWarehouseNotFound
	}
	return w, nil
}

func (m *mockInventoryWarehouseRepo) FindBySellerID(sellerID string) ([]*entity.Warehouse, error) {
	return nil, nil
}

func (m *mockInventoryWarehouseRepo) Update(warehouse *entity.Warehouse) error {
	m.warehouses[warehouse.ID] = warehouse
	return nil
}

type mockInventoryCategoryRepo struct {
	access map[string]*entity.StoreCategory
}

func (m *mockInventoryCategoryRepo) Save(sc *entity.StoreCategory) error {
	m.access[sc.StoreID+":"+sc.CategoryID] = sc
	return nil
}

func (m *mockInventoryCategoryRepo) FindByStoreIDAndCategoryID(storeID, categoryID string) (*entity.StoreCategory, error) {
	sc, ok := m.access[storeID+":"+categoryID]
	if !ok {
		return nil, nil
	}
	return sc, nil
}

func (m *mockInventoryCategoryRepo) Update(sc *entity.StoreCategory) error {
	m.access[sc.StoreID+":"+sc.CategoryID] = sc
	return nil
}

type mockInventoryItemRepo struct {
	items []*entity.Inventory
}

func (m *mockInventoryItemRepo) Save(inv *entity.Inventory) error {
	m.items = append(m.items, inv)
	return nil
}

func (m *mockInventoryItemRepo) FindByID(id string) (*entity.Inventory, error) {
	return nil, nil
}

func (m *mockInventoryItemRepo) FindByStoreID(storeID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockInventoryItemRepo) FindByWarehouseID(warehouseID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockInventoryItemRepo) FindByProductID(productID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockInventoryItemRepo) Update(inv *entity.Inventory) error {
	return nil
}

func newActiveStoreForInventory(storeID, userID string) *entity.Store {
	s, _ := entity.NewStore(storeID, userID, "Electronics Shop", nil, nil, nil)
	s.Approve()
	s.Events()
	return s
}

func newLinkedWarehouseForInventory(warehouseID, sellerID, storeID string) *entity.Warehouse {
	w, _ := entity.NewWarehouse(warehouseID, sellerID, "Warehouse", validAddress(), 1000, "public")
	w.LinkToStore(storeID, "primary")
	w.Events()
	return w
}

func newApprovedCategoryAccess(storeID, categoryID string) *entity.StoreCategory {
	sc, _ := entity.NewStoreCategory(storeID, categoryID)
	sc.Approve()
	sc.Events()
	return sc
}

func TestAddInventoryItem_Success(t *testing.T) {
	storeRepo := &mockInventoryStoreRepo{stores: map[string]*entity.Store{
		"store-1": newActiveStoreForInventory("store-1", "user-1"),
	}}
	warehouseRepo := &mockInventoryWarehouseRepo{warehouses: map[string]*entity.Warehouse{
		"wh-1": newLinkedWarehouseForInventory("wh-1", "user-1", "store-1"),
	}}
	categoryRepo := &mockInventoryCategoryRepo{access: map[string]*entity.StoreCategory{
		"store-1:cat-7": newApprovedCategoryAccess("store-1", "cat-7"),
	}}
	itemRepo := &mockInventoryItemRepo{}

	uc := usecase.NewAddInventoryItemUseCase(storeRepo, warehouseRepo, categoryRepo, itemRepo)

	input := usecase.AddInventoryItemInput{
		ID:          "inv-1",
		StoreID:     "store-1",
		WarehouseID: "wh-1",
		ProductID:   "prod-1",
		CategoryID:  "cat-7",
		BasePrice:   500000,
		InstantQty:  50,
		SaleModel:   "fixed",
		Condition:   "new",
		MinOrderQty: 1,
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.ID != "inv-1" {
		t.Errorf("expected inv-1, got %s", output.ID)
	}
	if output.StoreID != "store-1" {
		t.Errorf("expected store-1, got %s", output.StoreID)
	}
	if output.WarehouseID != "wh-1" {
		t.Errorf("expected wh-1, got %s", output.WarehouseID)
	}
	if output.ProductID != "prod-1" {
		t.Errorf("expected prod-1, got %s", output.ProductID)
	}
	if output.BasePrice != 500000 {
		t.Errorf("expected 500000, got %f", output.BasePrice)
	}
	if output.FinalPrice != 500000 {
		t.Errorf("expected final price equal to base price (500000), got %f", output.FinalPrice)
	}
	if output.InstantQty != 50 {
		t.Errorf("expected 50, got %d", output.InstantQty)
	}
	if output.VendorSaleStatus != string(entity.VendorSaleStatusActive) {
		t.Errorf("expected active vendor status, got %s", output.VendorSaleStatus)
	}
	if output.SystemSaleStatus != string(entity.SystemSaleStatusActive) {
		t.Errorf("expected active system status, got %s", output.SystemSaleStatus)
	}

	if len(itemRepo.items) != 1 {
		t.Fatalf("expected 1 saved inventory item, got %d", len(itemRepo.items))
	}
}

func TestAddInventoryItem_EventEmitted(t *testing.T) {
	storeRepo := &mockInventoryStoreRepo{stores: map[string]*entity.Store{
		"store-1": newActiveStoreForInventory("store-1", "user-1"),
	}}
	warehouseRepo := &mockInventoryWarehouseRepo{warehouses: map[string]*entity.Warehouse{
		"wh-1": newLinkedWarehouseForInventory("wh-1", "user-1", "store-1"),
	}}
	categoryRepo := &mockInventoryCategoryRepo{access: map[string]*entity.StoreCategory{
		"store-1:cat-7": newApprovedCategoryAccess("store-1", "cat-7"),
	}}
	itemRepo := &mockInventoryItemRepo{}

	uc := usecase.NewAddInventoryItemUseCase(storeRepo, warehouseRepo, categoryRepo, itemRepo)

	input := usecase.AddInventoryItemInput{
		ID:          "inv-1",
		StoreID:     "store-1",
		WarehouseID: "wh-1",
		ProductID:   "prod-1",
		CategoryID:  "cat-7",
		BasePrice:   500000,
		InstantQty:  50,
		SaleModel:   "fixed",
		Condition:   "new",
		MinOrderQty: 1,
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Event == nil {
		t.Fatal("expected a domain event, got nil")
	}
	e, ok := output.Event.(event.InventoryItemCreated)
	if !ok {
		t.Fatalf("expected InventoryItemCreated, got %T", output.Event)
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected inv-1, got %s", e.InventoryID)
	}
	if e.EventName() != "inventory.item_created" {
		t.Errorf("expected inventory.item_created, got %s", e.EventName())
	}
}

func TestAddInventoryItem_StoreNotFound(t *testing.T) {
	storeRepo := &mockInventoryStoreRepo{stores: map[string]*entity.Store{}}
	warehouseRepo := &mockInventoryWarehouseRepo{warehouses: map[string]*entity.Warehouse{
		"wh-1": newLinkedWarehouseForInventory("wh-1", "user-1", "store-1"),
	}}
	categoryRepo := &mockInventoryCategoryRepo{access: map[string]*entity.StoreCategory{
		"store-1:cat-7": newApprovedCategoryAccess("store-1", "cat-7"),
	}}
	itemRepo := &mockInventoryItemRepo{}

	uc := usecase.NewAddInventoryItemUseCase(storeRepo, warehouseRepo, categoryRepo, itemRepo)

	input := usecase.AddInventoryItemInput{
		ID:          "inv-1",
		StoreID:     "store-1",
		WarehouseID: "wh-1",
		ProductID:   "prod-1",
		CategoryID:  "cat-7",
		BasePrice:   500000,
		InstantQty:  50,
		SaleModel:   "fixed",
		Condition:   "new",
		MinOrderQty: 1,
	}

	_, err := uc.Execute(input)
	if err != entity.ErrStoreNotFound {
		t.Errorf("expected ErrStoreNotFound, got %v", err)
	}
}

func TestAddInventoryItem_StoreNotActive(t *testing.T) {
	s, _ := entity.NewStore("store-1", "user-1", "Electronics Shop", nil, nil, nil)
	s.Events()

	storeRepo := &mockInventoryStoreRepo{stores: map[string]*entity.Store{"store-1": s}}
	warehouseRepo := &mockInventoryWarehouseRepo{warehouses: map[string]*entity.Warehouse{
		"wh-1": newLinkedWarehouseForInventory("wh-1", "user-1", "store-1"),
	}}
	categoryRepo := &mockInventoryCategoryRepo{access: map[string]*entity.StoreCategory{
		"store-1:cat-7": newApprovedCategoryAccess("store-1", "cat-7"),
	}}
	itemRepo := &mockInventoryItemRepo{}

	uc := usecase.NewAddInventoryItemUseCase(storeRepo, warehouseRepo, categoryRepo, itemRepo)

	input := usecase.AddInventoryItemInput{
		ID:          "inv-1",
		StoreID:     "store-1",
		WarehouseID: "wh-1",
		ProductID:   "prod-1",
		CategoryID:  "cat-7",
		BasePrice:   500000,
		InstantQty:  50,
		SaleModel:   "fixed",
		Condition:   "new",
		MinOrderQty: 1,
	}

	_, err := uc.Execute(input)
	if err != usecase.ErrStoreNotActive {
		t.Errorf("expected ErrStoreNotActive, got %v", err)
	}
}

func TestAddInventoryItem_WarehouseNotFound(t *testing.T) {
	storeRepo := &mockInventoryStoreRepo{stores: map[string]*entity.Store{
		"store-1": newActiveStoreForInventory("store-1", "user-1"),
	}}
	warehouseRepo := &mockInventoryWarehouseRepo{warehouses: map[string]*entity.Warehouse{}}
	categoryRepo := &mockInventoryCategoryRepo{access: map[string]*entity.StoreCategory{
		"store-1:cat-7": newApprovedCategoryAccess("store-1", "cat-7"),
	}}
	itemRepo := &mockInventoryItemRepo{}

	uc := usecase.NewAddInventoryItemUseCase(storeRepo, warehouseRepo, categoryRepo, itemRepo)

	input := usecase.AddInventoryItemInput{
		ID:          "inv-1",
		StoreID:     "store-1",
		WarehouseID: "wh-1",
		ProductID:   "prod-1",
		CategoryID:  "cat-7",
		BasePrice:   500000,
		InstantQty:  50,
		SaleModel:   "fixed",
		Condition:   "new",
		MinOrderQty: 1,
	}

	_, err := uc.Execute(input)
	if err != entity.ErrWarehouseNotFound {
		t.Errorf("expected ErrWarehouseNotFound, got %v", err)
	}
}

func TestAddInventoryItem_WarehouseNotLinkedToStore(t *testing.T) {
	storeRepo := &mockInventoryStoreRepo{stores: map[string]*entity.Store{
		"store-1": newActiveStoreForInventory("store-1", "user-1"),
	}}
	warehouseRepo := &mockInventoryWarehouseRepo{warehouses: map[string]*entity.Warehouse{
		"wh-1": newLinkedWarehouseForInventory("wh-1", "user-1", "other-store"),
	}}
	categoryRepo := &mockInventoryCategoryRepo{access: map[string]*entity.StoreCategory{
		"store-1:cat-7": newApprovedCategoryAccess("store-1", "cat-7"),
	}}
	itemRepo := &mockInventoryItemRepo{}

	uc := usecase.NewAddInventoryItemUseCase(storeRepo, warehouseRepo, categoryRepo, itemRepo)

	input := usecase.AddInventoryItemInput{
		ID:          "inv-1",
		StoreID:     "store-1",
		WarehouseID: "wh-1",
		ProductID:   "prod-1",
		CategoryID:  "cat-7",
		BasePrice:   500000,
		InstantQty:  50,
		SaleModel:   "fixed",
		Condition:   "new",
		MinOrderQty: 1,
	}

	_, err := uc.Execute(input)
	if err != usecase.ErrWarehouseNotLinkedToStore {
		t.Errorf("expected ErrWarehouseNotLinkedToStore, got %v", err)
	}
}

func TestAddInventoryItem_CategoryAccessNotApproved(t *testing.T) {
	storeRepo := &mockInventoryStoreRepo{stores: map[string]*entity.Store{
		"store-1": newActiveStoreForInventory("store-1", "user-1"),
	}}
	warehouseRepo := &mockInventoryWarehouseRepo{warehouses: map[string]*entity.Warehouse{
		"wh-1": newLinkedWarehouseForInventory("wh-1", "user-1", "store-1"),
	}}
	categoryRepo := &mockInventoryCategoryRepo{access: map[string]*entity.StoreCategory{}}
	itemRepo := &mockInventoryItemRepo{}

	uc := usecase.NewAddInventoryItemUseCase(storeRepo, warehouseRepo, categoryRepo, itemRepo)

	input := usecase.AddInventoryItemInput{
		ID:          "inv-1",
		StoreID:     "store-1",
		WarehouseID: "wh-1",
		ProductID:   "prod-1",
		CategoryID:  "cat-7",
		BasePrice:   500000,
		InstantQty:  50,
		SaleModel:   "fixed",
		Condition:   "new",
		MinOrderQty: 1,
	}

	_, err := uc.Execute(input)
	if err != usecase.ErrCategoryAccessNotApproved {
		t.Errorf("expected ErrCategoryAccessNotApproved, got %v", err)
	}
}

func TestAddInventoryItem_CategoryAccessPending(t *testing.T) {
	storeRepo := &mockInventoryStoreRepo{stores: map[string]*entity.Store{
		"store-1": newActiveStoreForInventory("store-1", "user-1"),
	}}
	warehouseRepo := &mockInventoryWarehouseRepo{warehouses: map[string]*entity.Warehouse{
		"wh-1": newLinkedWarehouseForInventory("wh-1", "user-1", "store-1"),
	}}
	sc, _ := entity.NewStoreCategory("store-1", "cat-7")
	sc.Events()
	categoryRepo := &mockInventoryCategoryRepo{access: map[string]*entity.StoreCategory{
		"store-1:cat-7": sc,
	}}
	itemRepo := &mockInventoryItemRepo{}

	uc := usecase.NewAddInventoryItemUseCase(storeRepo, warehouseRepo, categoryRepo, itemRepo)

	input := usecase.AddInventoryItemInput{
		ID:          "inv-1",
		StoreID:     "store-1",
		WarehouseID: "wh-1",
		ProductID:   "prod-1",
		CategoryID:  "cat-7",
		BasePrice:   500000,
		InstantQty:  50,
		SaleModel:   "fixed",
		Condition:   "new",
		MinOrderQty: 1,
	}

	_, err := uc.Execute(input)
	if err != usecase.ErrCategoryAccessNotApproved {
		t.Errorf("expected ErrCategoryAccessNotApproved, got %v", err)
	}
}

func TestAddInventoryItem_InvalidInput(t *testing.T) {
	storeRepo := &mockInventoryStoreRepo{stores: map[string]*entity.Store{
		"store-1": newActiveStoreForInventory("store-1", "user-1"),
	}}
	warehouseRepo := &mockInventoryWarehouseRepo{warehouses: map[string]*entity.Warehouse{
		"wh-1": newLinkedWarehouseForInventory("wh-1", "user-1", "store-1"),
	}}
	categoryRepo := &mockInventoryCategoryRepo{access: map[string]*entity.StoreCategory{
		"store-1:cat-7": newApprovedCategoryAccess("store-1", "cat-7"),
	}}
	itemRepo := &mockInventoryItemRepo{}

	uc := usecase.NewAddInventoryItemUseCase(storeRepo, warehouseRepo, categoryRepo, itemRepo)

	tests := []struct {
		name  string
		input usecase.AddInventoryItemInput
	}{
		{"empty id", usecase.AddInventoryItemInput{ID: "", StoreID: "store-1", WarehouseID: "wh-1", ProductID: "prod-1", CategoryID: "cat-7", BasePrice: 100, InstantQty: 10, SaleModel: "fixed", Condition: "new", MinOrderQty: 1}},
		{"empty store id", usecase.AddInventoryItemInput{ID: "inv-1", StoreID: "", WarehouseID: "wh-1", ProductID: "prod-1", CategoryID: "cat-7", BasePrice: 100, InstantQty: 10, SaleModel: "fixed", Condition: "new", MinOrderQty: 1}},
		{"empty warehouse id", usecase.AddInventoryItemInput{ID: "inv-1", StoreID: "store-1", WarehouseID: "", ProductID: "prod-1", CategoryID: "cat-7", BasePrice: 100, InstantQty: 10, SaleModel: "fixed", Condition: "new", MinOrderQty: 1}},
		{"empty product id", usecase.AddInventoryItemInput{ID: "inv-1", StoreID: "store-1", WarehouseID: "wh-1", ProductID: "", CategoryID: "cat-7", BasePrice: 100, InstantQty: 10, SaleModel: "fixed", Condition: "new", MinOrderQty: 1}},
		{"zero base price", usecase.AddInventoryItemInput{ID: "inv-1", StoreID: "store-1", WarehouseID: "wh-1", ProductID: "prod-1", CategoryID: "cat-7", BasePrice: 0, InstantQty: 10, SaleModel: "fixed", Condition: "new", MinOrderQty: 1}},
		{"negative stock", usecase.AddInventoryItemInput{ID: "inv-1", StoreID: "store-1", WarehouseID: "wh-1", ProductID: "prod-1", CategoryID: "cat-7", BasePrice: 100, InstantQty: -1, SaleModel: "fixed", Condition: "new", MinOrderQty: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(tt.input)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
