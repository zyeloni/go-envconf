package config

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

// TestParseTag sprawdza funkcję parseTag
func TestParseTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected map[string]string
	}{
		{
			name: "Empty tag",
			tag:  "",
			expected: map[string]string{},
		},
		{
			name: "Single key-value pair",
			tag:  "env=TEST_VAR",
			expected: map[string]string{
				"env": "TEST_VAR",
			},
		},
		{
			name: "Multiple key-value pairs",
			tag:  "env=TEST_VAR,default=default value,required=true",
			expected: map[string]string{
				"env":      "TEST_VAR",
				"default":  "default value",
				"required": "true",
			},
		},
		{
			name: "Whitespace handling",
			tag:  " env = TEST_VAR , default = default value ",
			expected: map[string]string{
				"env":     "TEST_VAR",
				"default": "default value",
			},
		},
		{
			name: "Invalid format (no value)",
			tag:  "env",
			expected: map[string]string{},
		},
		{
			name: "Mixed valid and invalid",
			tag:  "env=TEST_VAR,invalid,default=value",
			expected: map[string]string{
				"env":     "TEST_VAR",
				"default": "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTag(tt.tag)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseTag() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSetFieldValue sprawdza funkcję setFieldValue
func TestSetFieldValue(t *testing.T) {
	// Test dla string
	t.Run("String", func(t *testing.T) {
		// Tworzymy strukturę z polem string
		type TestStruct struct {
			Field string
		}
		s := &TestStruct{}
		field := reflect.ValueOf(s).Elem().Field(0)
		
		err := setFieldValue(field, "test string", "Field")
		if err != nil {
			t.Errorf("setFieldValue() error = %v", err)
		}
		if s.Field != "test string" {
			t.Errorf("setFieldValue() = %v, want %v", s.Field, "test string")
		}
	})

	// Test dla int
	t.Run("Int", func(t *testing.T) {
		// Tworzymy strukturę z polem int
		type TestStruct struct {
			Field int
		}
		s := &TestStruct{}
		field := reflect.ValueOf(s).Elem().Field(0)
		
		err := setFieldValue(field, "42", "Field")
		if err != nil {
			t.Errorf("setFieldValue() error = %v", err)
		}
		if s.Field != 42 {
			t.Errorf("setFieldValue() = %v, want %v", s.Field, 42)
		}
	})

	// Test dla float
	t.Run("Float", func(t *testing.T) {
		// Tworzymy strukturę z polem float64
		type TestStruct struct {
			Field float64
		}
		s := &TestStruct{}
		field := reflect.ValueOf(s).Elem().Field(0)
		
		err := setFieldValue(field, "3.14", "Field")
		if err != nil {
			t.Errorf("setFieldValue() error = %v", err)
		}
		if s.Field != 3.14 {
			t.Errorf("setFieldValue() = %v, want %v", s.Field, 3.14)
		}
	})

	// Test dla bool
	t.Run("Bool", func(t *testing.T) {
		// Tworzymy strukturę z polem bool
		type TestStruct struct {
			Field bool
		}
		s := &TestStruct{}
		field := reflect.ValueOf(s).Elem().Field(0)
		
		err := setFieldValue(field, "true", "Field")
		if err != nil {
			t.Errorf("setFieldValue() error = %v", err)
		}
		if s.Field != true {
			t.Errorf("setFieldValue() = %v, want %v", s.Field, true)
		}
	})

	// Test dla time.Duration
	t.Run("Duration", func(t *testing.T) {
		// Tworzymy strukturę z polem time.Duration
		type TestStruct struct {
			Field time.Duration
		}
		s := &TestStruct{}
		field := reflect.ValueOf(s).Elem().Field(0)
		
		err := setFieldValue(field, "5s", "Field")
		if err != nil {
			t.Errorf("setFieldValue() error = %v", err)
		}
		if s.Field != 5*time.Second {
			t.Errorf("setFieldValue() = %v, want %v", s.Field, 5*time.Second)
		}
	})

	// Test dla time.Time
	t.Run("Time", func(t *testing.T) {
		// Tworzymy strukturę z polem time.Time
		type TestStruct struct {
			Field time.Time
		}
		s := &TestStruct{}
		field := reflect.ValueOf(s).Elem().Field(0)
		
		err := setFieldValue(field, "2023-01-02T15:04:05Z", "Field")
		if err != nil {
			t.Errorf("setFieldValue() error = %v", err)
		}
		expected, _ := time.Parse(time.RFC3339, "2023-01-02T15:04:05Z")
		if !s.Field.Equal(expected) {
			t.Errorf("setFieldValue() = %v, want %v", s.Field, expected)
		}
	})

	// Test dla nieprawidłowej wartości
	t.Run("Invalid value", func(t *testing.T) {
		// Tworzymy strukturę z polem int
		type TestStruct struct {
			Field int
		}
		s := &TestStruct{}
		field := reflect.ValueOf(s).Elem().Field(0)
		
		err := setFieldValue(field, "not an int", "Field")
		if err == nil {
			t.Errorf("setFieldValue() error = nil, want error")
		}
		
		var parseErr *ParseError
		if !errors.As(err, &parseErr) {
			t.Errorf("setFieldValue() error type = %T, want *ParseError", err)
		} else {
			// Sprawdź pola błędu
			if parseErr.FieldName != "Field" {
				t.Errorf("ParseError.FieldName = %v, want %v", parseErr.FieldName, "Field")
			}
			if parseErr.Value != "not an int" {
				t.Errorf("ParseError.Value = %v, want %v", parseErr.Value, "not an int")
			}
		}
	})

	// Test dla nieobsługiwanego typu
	t.Run("Unsupported type", func(t *testing.T) {
		// Tworzymy strukturę z polem chan int
		type TestStruct struct {
			Field chan int
		}
		s := &TestStruct{}
		field := reflect.ValueOf(s).Elem().Field(0)
		
		err := setFieldValue(field, "value", "Field")
		if err == nil {
			t.Errorf("setFieldValue() error = nil, want error")
		}
	})
}

// TestLoadStruct sprawdza funkcję LoadStruct
func TestLoadStruct(t *testing.T) {
	// Test dla prostej struktury
	t.Run("Simple struct", func(t *testing.T) {
		type Config struct {
			String string `config:"env=TEST_STRING,default=default"`
			Int    int    `config:"env=TEST_INT,default=42"`
		}
		
		var cfg Config
		err := LoadStruct(reflect.ValueOf(&cfg).Elem())
		if err != nil {
			t.Errorf("LoadStruct() error = %v", err)
		}
		
		if cfg.String != "default" {
			t.Errorf("LoadStruct() String = %v, want %v", cfg.String, "default")
		}
		
		if cfg.Int != 42 {
			t.Errorf("LoadStruct() Int = %v, want %v", cfg.Int, 42)
		}
	})
	
	// Test dla wymaganego pola
	t.Run("Required field", func(t *testing.T) {
		type Config struct {
			Required string `config:"env=TEST_REQUIRED,required=true"`
		}
		
		var cfg Config
		err := LoadStruct(reflect.ValueOf(&cfg).Elem())
		if err == nil {
			t.Errorf("LoadStruct() error = nil, want RequiredFieldError")
		}
		
		var reqErr *RequiredFieldError
		if !errors.As(err, &reqErr) {
			t.Errorf("LoadStruct() error type = %T, want *RequiredFieldError", err)
		} else {
			// Sprawdź pola błędu
			if reqErr.FieldName != "Required" {
				t.Errorf("RequiredFieldError.FieldName = %v, want %v", reqErr.FieldName, "Required")
			}
			if reqErr.EnvName != "TEST_REQUIRED" {
				t.Errorf("RequiredFieldError.EnvName = %v, want %v", reqErr.EnvName, "TEST_REQUIRED")
			}
		}
	})
}