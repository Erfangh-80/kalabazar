package entity

import "errors"

// Address represents a real-world physical location.
// It is a shared value object usable across multiple entity types
// such as Store, Warehouse, and Seller.
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
	Latitude   *float64
	Longitude  *float64
}

// Validate checks that the address represents a valid physical location.
func (a Address) Validate() error {
	if a.Street == "" {
		return errors.New("street cannot be empty")
	}
	if a.City == "" {
		return errors.New("city cannot be empty")
	}
	if a.Country == "" {
		return errors.New("country cannot be empty")
	}
	if (a.Latitude == nil) != (a.Longitude == nil) {
		return errors.New("latitude and longitude must be provided together")
	}
	return nil
}
