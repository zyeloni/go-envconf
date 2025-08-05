// Package envconfig dostarcza funkcjonalność do ładowania konfiguracji z zmiennych środowiskowych
// do struktur Go przy użyciu tagów struktury. Wspiera różne typy danych, wartości domyślne
// oraz zagnieżdżone struktury dla lepszej organizacji konfiguracji.
package envconfig

import (
	"reflect"
)

// Stałe używane do parsowania tagów struktury
const (
	Tag         = "envconfig" // Nazwa tagu używanego do konfiguracji
	EnvKey      = "env"       // Klucz określający nazwę zmiennej środowiskowej
	DefaultKey  = "default"   // Klucz określający wartość domyślną
	RequiredKey = "required"  // Klucz określający czy pole jest wymagane
)

// Load ładuje konfigurację z zmiennych środowiskowych do podanej struktury.
// Parametr config musi być wskaźnikiem do struktury, w przeciwnym razie zostanie zwrócony błąd.
// Funkcja przeszukuje wszystkie pola struktury i ustawia ich wartości na podstawie zmiennych środowiskowych
// lub wartości domyślnych określonych w tagach struktury.
// Jeśli pole jest oznaczone jako wymagane (required=true), a nie ma wartości, zwraca błąd.
func Load(config interface{}) error {
	configValue := reflect.ValueOf(config)
	// Sprawdzenie czy config jest wskaźnikiem do struktury
	if configValue.Kind() != reflect.Ptr || configValue.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}

	return LoadStruct(configValue.Elem())
}
