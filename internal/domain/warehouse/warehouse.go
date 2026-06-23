package warehouse

import "time"

type Warehouse struct {
	ID        int64
	Name      string
	Capacity  int
	CreatedAt time.Time
	events    []any
}

func NewWarehouse(name string, capacity int) (*Warehouse, error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	w := &Warehouse{
		Name:      name,
		Capacity:  capacity,
		CreatedAt: time.Now(),
	}

	w.events = append(w.events, WarehouseCreatedEvent{
		WarehouseID: w.ID,
		Name:        name,
	})

	return w, nil
}

func (w *Warehouse) Events() []any {
	return w.events
}

func (w *Warehouse) ClearEvents() {
	w.events = nil
}
