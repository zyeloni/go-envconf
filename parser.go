// Package config dostarcza funkcjonalność do ładowania konfiguracji z zmiennych środowiskowych
// do struktur Go przy użyciu tagów struktury.
package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// LoadStruct rekurencyjnie ładuje wartości ze zmiennych środowiskowych do pól struktury.
// Funkcja przechodzi przez wszystkie pola struktury i dla każdego pola z tagiem "config"
// próbuje załadować wartość z odpowiedniej zmiennej środowiskowej lub użyć wartości domyślnej.
// Jeśli pole jest oznaczone jako wymagane (required = true), a nie ma wartości, zwraca błąd.
func LoadStruct(structValue reflect.Value) error {
	structType := structValue.Type()

	// Iteracja przez wszystkie pola struktury
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Pomijamy pola, których nie można ustawić (np. nieeksportowane)
		if !field.CanSet() {
			continue
		}

		// Rekurencyjne przetwarzanie zagnieżdżonych struktur bez własnego tagu config
		if field.Kind() == reflect.Struct && fieldType.Tag.Get(Tag) == "" {
			if err := LoadStruct(field); err != nil {
				return err
			}
			continue
		}

		// Pobierz tag konfiguracji dla pola
		tag := fieldType.Tag.Get(Tag)
		if tag == "" {
			continue
		}

		// Parsowanie tagu do mapy klucz-wartość
		tagMap := parseTag(tag)

		// Ustalenie nazwy zmiennej środowiskowej
		envName, ok := tagMap[EnvKey]
		if !ok {
			// Jeśli nie określono nazwy zmiennej, użyj nazwy pola w górnym rejestrze
			envName = strings.ToUpper(fieldType.Name)
		}

		// Pobierz wartość z zmiennej środowiskowej
		envValue := os.Getenv(envName)

		// Jeśli zmienna środowiskowa nie jest ustawiona, użyj wartości domyślnej
		if envValue == "" {
			defaultValue, ok := tagMap[DefaultKey]
			if ok {
				envValue = defaultValue
			} else {
				// Sprawdź czy pole jest wymagane
				if required, ok := tagMap[RequiredKey]; ok && required == "true" {
					return &RequiredFieldError{
						FieldName: fieldType.Name,
						EnvName:   envName,
					}
				}
				continue
			}
		}

		// Ustaw wartość pola na podstawie wartości zmiennej środowiskowej
		if err := setFieldValue(field, envValue, fieldType.Name); err != nil {
			return err
		}
	}

	return nil
}

// parseTag parsuje tag struktury w formacie "klucz1=wartość1,klucz2=wartość2"
// i zwraca mapę par klucz-wartość.
func parseTag(tag string) map[string]string {
	result := make(map[string]string)
	// Podział tagu na części oddzielone przecinkami
	parts := strings.Split(tag, ",")

	// Przetwarzanie każdej części jako pary klucz=wartość
	for _, part := range parts {
		// Podział na klucz i wartość przy pierwszym znaku "="
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			// Dodanie pary do mapy wynikowej, usuwając białe znaki
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	return result
}

// setFieldValue ustawia wartość pola struktury na podstawie wartości tekstowej.
// Funkcja obsługuje różne typy danych, w tym string, int, uint, float, bool, time.Time i time.Duration.
// Dla nieobsługiwanych typów zwraca błąd.
func setFieldValue(field reflect.Value, value string, fieldName string) error {
	// Obsługa specjalnego typu time.Time
	if field.Type() == reflect.TypeOf(time.Time{}) {
		timeValue, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return &ParseError{
				FieldName: fieldName,
				FieldType: "time.Time",
				Value:     value,
				Err:       err,
			}
		}
		field.Set(reflect.ValueOf(timeValue))
		return nil
	}

	// Obsługa specjalnego typu time.Duration
	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		durationValue, err := time.ParseDuration(value)
		if err != nil {
			return &ParseError{
				FieldName: fieldName,
				FieldType: "time.Duration",
				Value:     value,
				Err:       err,
			}
		}
		field.Set(reflect.ValueOf(durationValue))
		return nil
	}

	// Obsługa standardowych typów Go na podstawie rodzaju pola
	switch field.Kind() {
	case reflect.String:
		// Dla stringów bezpośrednio ustawiamy wartość
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Konwersja string -> int
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return &ParseError{
				FieldName: fieldName,
				FieldType: field.Kind().String(),
				Value:     value,
				Err:       err,
			}
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Konwersja string -> uint
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return &ParseError{
				FieldName: fieldName,
				FieldType: field.Kind().String(),
				Value:     value,
				Err:       err,
			}
		}
		field.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		// Konwersja string -> float
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return &ParseError{
				FieldName: fieldName,
				FieldType: field.Kind().String(),
				Value:     value,
				Err:       err,
			}
		}
		field.SetFloat(floatValue)
	case reflect.Bool:
		// Konwersja string -> bool
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return &ParseError{
				FieldName: fieldName,
				FieldType: field.Kind().String(),
				Value:     value,
				Err:       err,
			}
		}
		field.SetBool(boolValue)
	case reflect.Struct:
		// Rekurencyjne przetwarzanie zagnieżdżonych struktur
		return LoadStruct(field)
	default:
		// Zwróć błąd dla nieobsługiwanych typów
		return fmt.Errorf("%w: %s", ErrUnsupportedFieldType, field.Kind().String())
	}
	return nil
}
