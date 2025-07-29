// Package config dostarcza funkcjonalność do ładowania konfiguracji z zmiennych środowiskowych
// do struktur Go przy użyciu tagów struktury.
package config

import (
	"errors"
	"fmt"
)

// Podstawowe błędy zwracane przez pakiet
var (
	// ErrNotStruct zwracany gdy konfiguracja nie jest wskaźnikiem do struktury
	ErrNotStruct = errors.New("config must be a pointer to a struct")

	// ErrUnsupportedFieldType zwracany gdy typ pola nie jest obsługiwany
	ErrUnsupportedFieldType = errors.New("unsupported field type")

	// ErrMissingRequired zwracany gdy wymagane pole nie ma wartości
	ErrMissingRequired = errors.New("missing required field")
)

// RequiredFieldError reprezentuje błąd brakującego wymaganego pola
type RequiredFieldError struct {
	FieldName string
	EnvName   string
}

// Error implementuje interfejs error
func (e *RequiredFieldError) Error() string {
	return fmt.Sprintf("%s: field '%s' is required but no value was provided (env: %s)",
		ErrMissingRequired.Error(), e.FieldName, e.EnvName)
}

// ParseError reprezentuje błąd podczas parsowania wartości
type ParseError struct {
	FieldName string
	FieldType string
	Value     string
	Err       error
}

// Error implementuje interfejs error
func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse value '%s' as %s for field '%s': %v",
		e.Value, e.FieldType, e.FieldName, e.Err)
}

// Unwrap implementuje interfejs errors.Unwrap
func (e *ParseError) Unwrap() error {
	return e.Err
}
