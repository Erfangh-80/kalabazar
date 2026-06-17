package entity_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
)

func TestAddress_Validate_Success(t *testing.T) {
	addr := entity.Address{
		Street: "123 Main St", City: "Tehran", Country: "Iran",
	}
	if err := addr.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAddress_Validate_WithOptionalFields(t *testing.T) {
	lat := 35.6892
	lng := 51.3890
	addr := entity.Address{
		Street: "123 Main St", City: "Tehran",
		State: "Tehran Province", PostalCode: "1234567890",
		Country: "Iran", Latitude: &lat, Longitude: &lng,
	}
	if err := addr.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAddress_Validate_EmptyStreet(t *testing.T) {
	addr := entity.Address{Street: "", City: "Tehran", Country: "Iran"}
	if err := addr.Validate(); err == nil {
		t.Fatal("expected error for empty street")
	}
}

func TestAddress_Validate_EmptyCity(t *testing.T) {
	addr := entity.Address{Street: "123 Main St", City: "", Country: "Iran"}
	if err := addr.Validate(); err == nil {
		t.Fatal("expected error for empty city")
	}
}

func TestAddress_Validate_EmptyCountry(t *testing.T) {
	addr := entity.Address{Street: "123 Main St", City: "Tehran", Country: ""}
	if err := addr.Validate(); err == nil {
		t.Fatal("expected error for empty country")
	}
}

func TestAddress_Validate_MissingLatitude(t *testing.T) {
	lng := 51.3890
	addr := entity.Address{
		Street: "123 Main St", City: "Tehran",
		Country: "Iran", Longitude: &lng,
	}
	if err := addr.Validate(); err == nil {
		t.Fatal("expected error when longitude provided without latitude")
	}
}

func TestAddress_Validate_MissingLongitude(t *testing.T) {
	lat := 35.6892
	addr := entity.Address{
		Street: "123 Main St", City: "Tehran",
		Country: "Iran", Latitude: &lat,
	}
	if err := addr.Validate(); err == nil {
		t.Fatal("expected error when latitude provided without longitude")
	}
}

func TestAddress_Validate_AllEmpty(t *testing.T) {
	addr := entity.Address{}
	if err := addr.Validate(); err == nil {
		t.Fatal("expected error for empty address")
	}
}
