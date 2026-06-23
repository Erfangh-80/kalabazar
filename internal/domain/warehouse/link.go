package warehouse

import "time"

type StoreWarehouseLink struct {
	ID          int64
	StoreID     int64
	WarehouseID int64
	LinkType    string
	CreatedAt   time.Time
	events      []any
}

func NewStoreWarehouseLink(storeID, warehouseID int64, linkType string) (*StoreWarehouseLink, error) {
	link := &StoreWarehouseLink{
		StoreID:     storeID,
		WarehouseID: warehouseID,
		LinkType:    linkType,
		CreatedAt:   time.Now(),
	}

	link.events = append(link.events, WarehouseLinkedToStoreEvent{
		WarehouseID: warehouseID,
		StoreID:     storeID,
		LinkType:    linkType,
	})

	return link, nil
}

func (l *StoreWarehouseLink) Events() []any {
	return l.events
}

func (l *StoreWarehouseLink) ClearEvents() {
	l.events = nil
}
