package envconfig

import (
	"errors"
	"os"
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
			name:     "Empty tag",
			tag:      "",
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
			name:     "Invalid format (no value)",
			tag:      "env",
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
		t.Run(
			tt.name, func(t *testing.T) {
				result := parseTag(tt.tag)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("parseTag() = %v, want %v", result, tt.expected)
				}
			},
		)
	}
}

// TestSetFieldValue sprawdza funkcję setFieldValue
func TestSetFieldValue(t *testing.T) {
	// Test dla string
	t.Run(
		"String", func(t *testing.T) {
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
		},
	)

	// Test dla int
	t.Run(
		"Int", func(t *testing.T) {
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
		},
	)

	// Test dla float
	t.Run(
		"Float", func(t *testing.T) {
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
		},
	)

	t.Run(
		"Uint", func(t *testing.T) {
			// Tworzymy strukturę z polem float64
			type TestStruct struct {
				Field uint
			}
			s := &TestStruct{}
			field := reflect.ValueOf(s).Elem().Field(0)

			err := setFieldValue(field, "3", "Field")
			if err != nil {
				t.Errorf("setFieldValue() error = %v", err)
			}
			if s.Field != 3 {
				t.Errorf("setFieldValue() = %v, want %v", s.Field, 3)
			}
		},
	)

	// Test dla bool
	t.Run(
		"Bool", func(t *testing.T) {
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
		},
	)

	// Test dla time.Duration
	t.Run(
		"Duration", func(t *testing.T) {
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
		},
	)

	// Test dla time.Time
	t.Run(
		"Time", func(t *testing.T) {
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
		},
	)

	// Test dla nieprawidłowej wartości
	t.Run(
		"Invalid value", func(t *testing.T) {
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
		},
	)

	// Test dla time.Time z nieprawidłowym formatem
	t.Run(
		"Invalid time format", func(t *testing.T) {
			// Tworzymy strukturę z polem time.Time
			type TestStruct struct {
				Field time.Time
			}
			s := &TestStruct{}
			field := reflect.ValueOf(s).Elem().Field(0)
			err := setFieldValue(field, "not a time", "Field")
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
				if parseErr.Value != "not a time" {
					t.Errorf("ParseError.Value = %v, want %v", parseErr.Value, "not a time")
				}
				if parseErr.FieldType != "time.Time" {
					t.Errorf("ParseError.FieldType = %v, want %v", parseErr.FieldType, "time.Time")
				}
			}
		},
	)

	// Test dla time.Duration z nieprawidłowym formatem
	t.Run(
		"Invalid duration format", func(t *testing.T) {
			// Tworzymy strukturę z polem time.Duration
			type TestStruct struct {
				Field time.Duration
			}
			s := &TestStruct{}
			field := reflect.ValueOf(s).Elem().Field(0)
			err := setFieldValue(field, "not a duration", "Field")
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
				if parseErr.Value != "not a duration" {
					t.Errorf("ParseError.Value = %v, want %v", parseErr.Value, "not a duration")
				}
				if parseErr.FieldType != "time.Duration" {
					t.Errorf("ParseError.FieldType = %v, want %v", parseErr.FieldType, "time.Duration")
				}
			}
		},
	)

	// Test dla bool z nieprawidłową wartością
	t.Run(
		"Invalid bool value", func(t *testing.T) {
			// Tworzymy strukturę z polem bool
			type TestStruct struct {
				Field bool
			}
			s := &TestStruct{}
			field := reflect.ValueOf(s).Elem().Field(0)
			err := setFieldValue(field, "not a bool", "Field")
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
				if parseErr.Value != "not a bool" {
					t.Errorf("ParseError.Value = %v, want %v", parseErr.Value, "not a bool")
				}
				if parseErr.FieldType != "bool" {
					t.Errorf("ParseError.FieldType = %v, want %v", parseErr.FieldType, "bool")
				}
			}
		},
	)

	// Test dla Uint z nieprawidłową wartością
	t.Run(
		"Invalid uint value", func(t *testing.T) {
			// Tworzymy strukturę z polem uint
			type TestStruct struct {
				Field uint
			}
			s := &TestStruct{}
			field := reflect.ValueOf(s).Elem().Field(0)
			err := setFieldValue(field, "not a uint", "Field")
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
				if parseErr.Value != "not a uint" {
					t.Errorf("ParseError.Value = %v, want %v", parseErr.Value, "not a uint")
				}
				if parseErr.FieldType != "uint" {
					t.Errorf("ParseError.FieldType = %v, want %v", parseErr.FieldType, "uint")
				}
			}
		},
	)

	// Test dla Float32 z nieprawidłową wartością
	t.Run(
		"Invalid float32 value", func(t *testing.T) {
			// Tworzymy strukturę z polem float32
			type TestStruct struct {
				Field float32
			}
			s := &TestStruct{}
			field := reflect.ValueOf(s).Elem().Field(0)
			err := setFieldValue(field, "not a float32", "Field")
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
				if parseErr.Value != "not a float32" {
					t.Errorf("ParseError.Value = %v, want %v", parseErr.Value, "not a float32")
				}
				if parseErr.FieldType != "float32" {
					t.Errorf("ParseError.FieldType = %v, want %v", parseErr.FieldType, "float32")
				}
			}
		},
	)

	// Test dla nieobsługiwanego typu
	t.Run(
		"Unsupported type", func(t *testing.T) {
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
		},
	)

}

// TestLoadStruct sprawdza funkcję LoadStruct
func TestLoadStruct(t *testing.T) {
	// Test dla prostej struktury
	t.Run(
		"Simple struct", func(t *testing.T) {
			type Config struct {
				String string `envconfig:"env=TEST_STRING,default=default"`
				Int    int    `envconfig:"env=TEST_INT,default=42"`
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
		},
	)

	// Test dla wymaganego pola
	t.Run(
		"Required field", func(t *testing.T) {
			type Config struct {
				Required string `envconfig:"env=TEST_REQUIRED,required=true"`
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
		},
	)

	// Test dla nieeksportowane pola
	t.Run(
		"Unexported field", func(t *testing.T) {
			type Config struct {
				unexportedField string `envconfig:"env=TEST_UNEXPORTED"`
			}

			var cfg Config
			err := LoadStruct(reflect.ValueOf(&cfg).Elem())
			if err != nil {
				t.Errorf("LoadStruct() error = %v", err)
			}

			if cfg.unexportedField != "" {
				t.Errorf("LoadStruct() unexportedField = %v, want empty string", cfg.unexportedField)
			}
		},
	)

	// Test dla struktury bez tagów
	t.Run(
		"Struct without tags", func(t *testing.T) {
			// Ustawiamy zmienną środowiskową dla testu
			os.Setenv("STRING", "test_value")
			defer os.Unsetenv("STRING")

			type Config struct {
				String string // Brak tagu
				Int    int    // Brak tagu
			}

			var cfg Config
			err := LoadStruct(reflect.ValueOf(&cfg).Elem())
			if err != nil {
				t.Errorf("LoadStruct() error = %v", err)
			}

			// Sprawdzamy czy wartość została poprawnie załadowana z zmiennej środowiskowej
			if cfg.String != "test_value" {
				t.Errorf("LoadStruct() String = %v, want %v", cfg.String, "test_value")
			}

			// Int powinien pozostać z wartością domyślną (0), ponieważ nie ma zmiennej środowiskowej INT
			if cfg.Int != 0 {
				t.Errorf("LoadStruct() Int = %v, want %v", cfg.Int, 0)
			}
		},
	)

	// Test dla zagnieżdżonej struktury bez tagów
	t.Run(
		"Nested struct without tags", func(t *testing.T) {
			// Ustawiamy zmienną środowiskową dla testu
			os.Setenv("NESTEDSTRING", "nested_value")
			defer os.Unsetenv("NESTEDSTRING")

			type NestedConfig struct {
				NestedString string // Brak tagu
			}

			type Config struct {
				Nested NestedConfig // Brak tagu
			}

			var cfg Config
			err := LoadStruct(reflect.ValueOf(&cfg).Elem())
			if err != nil {
				t.Errorf("LoadStruct() error = %v", err)
			}

			// Sprawdzamy czy wartość została poprawnie załadowana z zmiennej środowiskowej
			if cfg.Nested.NestedString != "nested_value" {
				t.Errorf("LoadStruct() Nested.NestedString = %v, want %v", cfg.Nested.NestedString, "nested_value")
			}
		},
	)
}
