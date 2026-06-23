package seller

import "unicode/utf8"

const (
	minNameLength = 2
	maxNameLength = 30
	phoneLength   = 11
)

func ValidateSellerName(name string) error {
	if name == "" {
		return ErrInvalidSellerName
	}
	if utf8.RuneCountInString(name) < minNameLength {
		return ErrInvalidSellerName
	}
	if utf8.RuneCountInString(name) > maxNameLength {
		return ErrInvalidSellerName
	}
	return nil
}

func ValidatePhone(phone string) error {
	if phone == "" {
		return ErrInvalidPhone
	}
	if len(phone) != phoneLength {
		return ErrInvalidPhone
	}
	for _, c := range phone {
		if c < '0' || c > '9' {
			return ErrInvalidPhone
		}
	}
	return nil
}
