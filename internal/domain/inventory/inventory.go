package inventory

import "time"

type Inventory struct {
	ID             int64
	ProductID      int64
	WarehouseID    int64
	BasePrice      int64
	FinalPrice     int64
	AvailableStock int
	ReservedStock  int
	StockOut       int
	CreatedAt      time.Time
	UpdatedAt      time.Time

	events []any
}

func NewInventory(productID, warehouseID int64, basePrice int64, stock int) (*Inventory, error) {
	if basePrice <= 0 {
		return nil, ErrInvalidPrice
	}
	if stock < 0 {
		return nil, ErrInvalidStock
	}

	now := time.Now()
	inv := &Inventory{
		ProductID:      productID,
		WarehouseID:    warehouseID,
		BasePrice:      basePrice,
		FinalPrice:     basePrice,
		AvailableStock: stock,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	inv.events = append(inv.events, InventoryCreatedEvent{
		ProductID:      productID,
		AvailableStock: stock,
		FinalPrice:     basePrice,
	})
	inv.events = append(inv.events, StockInEvent{
		Quantity: stock,
	})

	return inv, nil
}

func (inv *Inventory) ApplyDiscount(percentage int) {
	oldPrice := inv.FinalPrice
	discount := inv.BasePrice * int64(percentage) / 100
	inv.FinalPrice = inv.BasePrice - discount
	inv.UpdatedAt = time.Now()

	inv.events = append(inv.events, PriceUpdatedEvent{
		OldPrice: oldPrice,
		NewPrice: inv.FinalPrice,
	})
}

func (inv *Inventory) ReserveStock(qty int) error {
	if qty < 0 {
		return ErrInvalidQuantity
	}
	if qty > inv.AvailableStock {
		return ErrInsufficientStock
	}

	inv.AvailableStock -= qty
	inv.ReservedStock += qty
	inv.UpdatedAt = time.Now()

	inv.events = append(inv.events, ReservedEvent{
		Quantity: qty,
	})

	return nil
}

func (inv *Inventory) FinalizeSale(qty int) error {
	if qty < 0 {
		return ErrInvalidQuantity
	}
	if qty > inv.ReservedStock {
		return ErrInsufficientReservedStock
	}

	inv.ReservedStock -= qty
	inv.StockOut += qty
	inv.UpdatedAt = time.Now()

	inv.events = append(inv.events, StockOutEvent{
		Quantity: qty,
	})

	return nil
}

func (inv *Inventory) ResetPrice() {
	oldPrice := inv.FinalPrice
	inv.FinalPrice = inv.BasePrice
	inv.UpdatedAt = time.Now()

	inv.events = append(inv.events, PriceUpdatedEvent{
		OldPrice: oldPrice,
		NewPrice: inv.FinalPrice,
	})
}

func (inv *Inventory) RestoreStock(qty int) {
	if qty <= 0 {
		return
	}

	inv.AvailableStock += qty
	inv.ReservedStock -= qty
	inv.UpdatedAt = time.Now()
}

func (inv *Inventory) Events() []any {
	events := inv.events
	inv.events = nil
	return events
}
