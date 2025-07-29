package envconfig

import (
	"errors"
	"testing"
)

func TestRequiredFieldError_Error(t *testing.T) {
	err := &RequiredFieldError{
		FieldName: "DatabaseURL",
		EnvName:   "DATABASE_URL",
	}

	expected := "missing required field: field 'DatabaseURL' is required but no value was provided (env: DATABASE_URL)"
	if err.Error() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, err.Error())
	}
}

func TestParseError_Error(t *testing.T) {
	innerErr := errors.New("invalid syntax")
	err := &ParseError{
		FieldName: "Port",
		FieldType: "int",
		Value:     "abc",
		Err:       innerErr,
	}

	expected := "failed to parse value 'abc' as int for field 'Port': invalid syntax"
	if err.Error() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, err.Error())
	}
}

func TestParseError_Unwrap(t *testing.T) {
	innerErr := errors.New("invalid syntax")
	err := &ParseError{Err: innerErr}

	if !errors.Is(err, innerErr) {
		t.Error("Unwrap() did not return the expected inner error")
	}
}

func TestErrorsConstants(t *testing.T) {
	if ErrNotStruct.Error() != "config must be a pointer to a struct" {
		t.Error("ErrNotStruct does not match expected message")
	}
	if ErrUnsupportedFieldType.Error() != "unsupported field type" {
		t.Error("ErrUnsupportedFieldType does not match expected message")
	}
	if ErrMissingRequired.Error() != "missing required field" {
		t.Error("ErrMissingRequired does not match expected message")
	}
}
