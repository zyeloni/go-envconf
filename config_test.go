package envconfig

import (
	"errors"
	"os"
	"testing"
	"time"
)

// TestLoad_BasicTypes sprawdza ładowanie podstawowych typów danych
func TestLoad_BasicTypes(t *testing.T) {
	// Ustawienie zmiennych środowiskowych dla testu
	os.Setenv("TEST_STRING", "test value")
	os.Setenv("TEST_INT", "42")
	os.Setenv("TEST_FLOAT", "3.14")
	os.Setenv("TEST_BOOL", "true")
	os.Setenv("TEST_DURATION", "5s")
	os.Setenv("TEST_TIME", "2023-01-02T15:04:05Z")

	// Struktura testowa
	type Config struct {
		String   string        `envconfig:"env=TEST_STRING"`
		Int      int           `envconfig:"env=TEST_INT"`
		Float    float64       `envconfig:"env=TEST_FLOAT"`
		Bool     bool          `envconfig:"env=TEST_BOOL"`
		Duration time.Duration `envconfig:"env=TEST_DURATION"`
		Time     time.Time     `envconfig:"env=TEST_TIME"`
	}

	var cfg Config
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Sprawdzenie czy wartości zostały poprawnie załadowane
	if cfg.String != "test value" {
		t.Errorf("String = %v, want %v", cfg.String, "test value")
	}
	if cfg.Int != 42 {
		t.Errorf("Int = %v, want %v", cfg.Int, 42)
	}
	if cfg.Float != 3.14 {
		t.Errorf("Float = %v, want %v", cfg.Float, 3.14)
	}
	if cfg.Bool != true {
		t.Errorf("Bool = %v, want %v", cfg.Bool, true)
	}
	if cfg.Duration != 5*time.Second {
		t.Errorf("Duration = %v, want %v", cfg.Duration, 5*time.Second)
	}
	expectedTime, _ := time.Parse(time.RFC3339, "2023-01-02T15:04:05Z")
	if !cfg.Time.Equal(expectedTime) {
		t.Errorf("Time = %v, want %v", cfg.Time, expectedTime)
	}

	// Czyszczenie zmiennych środowiskowych
	os.Unsetenv("TEST_STRING")
	os.Unsetenv("TEST_INT")
	os.Unsetenv("TEST_FLOAT")
	os.Unsetenv("TEST_BOOL")
	os.Unsetenv("TEST_DURATION")
	os.Unsetenv("TEST_TIME")
}

// TestLoad_DefaultValues sprawdza użycie wartości domyślnych
func TestLoad_DefaultValues(t *testing.T) {
	// Struktura testowa z wartościami domyślnymi
	type Config struct {
		String string  `envconfig:"default=default value"`
		Int    int     `envconfig:"default=123"`
		Float  float64 `envconfig:"default=2.71"`
		Bool   bool    `envconfig:"default=true"`
	}

	var cfg Config
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Sprawdzenie czy wartości domyślne zostały poprawnie załadowane
	if cfg.String != "default value" {
		t.Errorf("String = %v, want %v", cfg.String, "default value")
	}
	if cfg.Int != 123 {
		t.Errorf("Int = %v, want %v", cfg.Int, 123)
	}
	if cfg.Float != 2.71 {
		t.Errorf("Float = %v, want %v", cfg.Float, 2.71)
	}
	if cfg.Bool != true {
		t.Errorf("Bool = %v, want %v", cfg.Bool, true)
	}
}

// TestLoad_RequiredFields sprawdza walidację wymaganych pól
func TestLoad_RequiredFields(t *testing.T) {
	// Struktura testowa z wymaganym polem
	type Config struct {
		Required string `envconfig:"env=TEST_REQUIRED,required=true"`
		Optional string `envconfig:"env=TEST_OPTIONAL"`
	}

	// Test 1: Brak wymaganego pola
	var cfg1 Config
	err := Load(&cfg1)
	if err == nil {
		t.Fatalf("Load() error = nil, want RequiredFieldError")
	}

	var reqErr *RequiredFieldError
	if !errors.As(err, &reqErr) {
		t.Fatalf("Load() error type = %T, want *RequiredFieldError", err)
	}

	if reqErr.FieldName != "Required" {
		t.Errorf("RequiredFieldError.FieldName = %v, want %v", reqErr.FieldName, "Required")
	}

	if reqErr.EnvName != "TEST_REQUIRED" {
		t.Errorf("RequiredFieldError.EnvName = %v, want %v", reqErr.EnvName, "TEST_REQUIRED")
	}

	// Test 2: Ustawienie wymaganego pola
	os.Setenv("TEST_REQUIRED", "required value")
	var cfg2 Config
	err = Load(&cfg2)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg2.Required != "required value" {
		t.Errorf("Required = %v, want %v", cfg2.Required, "required value")
	}

	// Czyszczenie zmiennych środowiskowych
	os.Unsetenv("TEST_REQUIRED")
}

// TestLoad_NestedStructs sprawdza ładowanie zagnieżdżonych struktur
func TestLoad_NestedStructs(t *testing.T) {
	// Ustawienie zmiennych środowiskowych dla testu
	os.Setenv("PARENT", "parent value")
	os.Setenv("CHILD", "child value")
	os.Setenv("GRANDCHILD", "grandchild value")

	// Struktura testowa z zagnieżdżonymi strukturami
	type GrandChild struct {
		Value string `envconfig:"env=GRANDCHILD"`
	}

	type Child struct {
		Value      string `envconfig:"env=CHILD"`
		GrandChild GrandChild
	}

	type Parent struct {
		Value string `envconfig:"env=PARENT"`
		Child Child
	}

	var cfg Parent
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Sprawdzenie czy wartości zostały poprawnie załadowane
	if cfg.Value != "parent value" {
		t.Errorf("Parent.Value = %v, want %v", cfg.Value, "parent value")
	}
	if cfg.Child.Value != "child value" {
		t.Errorf("Parent.Child.Value = %v, want %v", cfg.Child.Value, "child value")
	}
	if cfg.Child.GrandChild.Value != "grandchild value" {
		t.Errorf("Parent.Child.GrandChild.Value = %v, want %v", cfg.Child.GrandChild.Value, "grandchild value")
	}

	// Czyszczenie zmiennych środowiskowych
	os.Unsetenv("PARENT")
	os.Unsetenv("CHILD")
	os.Unsetenv("GRANDCHILD")
}

// TestLoad_InvalidInput sprawdza obsługę nieprawidłowych danych wejściowych
func TestLoad_InvalidInput(t *testing.T) {
	// Test 1: Nieprawidłowy typ konfiguracji (nie wskaźnik)
	var cfg1 struct{}
	err := Load(cfg1)
	if !errors.Is(err, ErrNotStruct) {
		t.Errorf("Load() error = %v, want %v", err, ErrNotStruct)
	}

	// Test 2: Nieprawidłowy typ konfiguracji (wskaźnik, ale nie do struktury)
	var i int
	err = Load(&i)
	if !errors.Is(err, ErrNotStruct) {
		t.Errorf("Load() error = %v, want %v", err, ErrNotStruct)
	}

	// Test 3: Nieprawidłowa wartość dla typu
	os.Setenv("TEST_INT", "not an int")
	type Config struct {
		Int int `envconfig:"env=TEST_INT"`
	}
	var cfg3 Config
	err = Load(&cfg3)
	if err == nil {
		t.Fatalf("Load() error = nil, want ParseError")
	}

	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("Load() error type = %T, want *ParseError", err)
	}

	// Czyszczenie zmiennych środowiskowych
	os.Unsetenv("TEST_INT")
}
