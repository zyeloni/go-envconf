# Config Manager

[![Test Status](https://github.com/zyeloni/go-envconf/actions/workflows/go.yml/badge.svg)](https://github.com/zyeloni/go-envconf/actions/workflows/go.yml)
[![Code Coverage](https://codecov.io/gh/zyeloni/go-envconf/branch/main/graph/badge.svg)](https://codecov.io/gh/zyeloni/go-envconf)
[![Go Reference](https://pkg.go.dev/badge/github.com/zyeloni/go-envconf.svg)](https://pkg.go.dev/github.com/zyeloni/go-envconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/zyeloni/go-envconf)](https://goreportcard.com/report/github.com/zyeloni/go-envconf)
[![Go Version](https://img.shields.io/github/go-mod/go-version/zyeloni/go-envconf)](https://github.com/zyeloni/go-envconf)
[![License](https://img.shields.io/github/license/zyeloni/go-envconf)](https://github.com/zyeloni/go-envconf/blob/main/LICENSE)

Prosta biblioteka Go do ładowania konfiguracji ze zmiennych środowiskowych przy użyciu tagów struktury.

## Funkcje

- Ładowanie konfiguracji ze zmiennych środowiskowych do struktur Go
- Definiowanie wartości domyślnych dla pól konfiguracyjnych
- Oznaczanie pól jako wymagane, aby zapewnić ich wartości
- Dostosowywanie nazw zmiennych środowiskowych za pomocą tagów struktury
- Obsługa różnych typów danych (string, int, uint, float, bool, time.Time, time.Duration)
- Obsługa zagnieżdżonych struktur dla lepszej organizacji konfiguracji
- Szczegółowe raportowanie błędów walidacji i parsowania
- Proste i łatwe w użyciu API

## Instalacja

```bash
go get github.com/zyeloni/go-envconf
```

Lub po prostu skopiuj pakiet `envconfig` do swojego projektu.

## CI/CD

Projekt korzysta z GitHub Actions do automatycznego testowania i weryfikacji jakości kodu. Workflow zawiera następujące etapy:

1. **Lint** - sprawdzanie jakości kodu za pomocą `golint` i `go vet`
2. **Test** - uruchamianie testów jednostkowych z pomiarem pokrycia kodu

Możesz zobaczyć status testów na odznace na górze tego README. Aby uruchomić testy lokalnie:

```bash
go test -v ./...
```

Aby uruchomić linting lokalnie:

```bash
go vet ./...
golint -set_exit_status ./...
```

## Użycie

### Podstawowe użycie

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/zyeloni/go-envconf"
)

// Zdefiniuj strukturę konfiguracyjną z tagami struktury
type AppConfig struct {
    ServerPort int    `config:"env=SERVER_PORT,default=8080"`
    ServerHost string `config:"env=SERVER_HOST,default=localhost"`
    Debug      bool   `config:"env=DEBUG,default=false"`
}

func main() {
    // Utwórz nową konfigurację z wartościami domyślnymi
    cfg := &AppConfig{}
    
    // Załaduj konfigurację ze zmiennych środowiskowych
    if err := envconfig.Load(cfg); err != nil {
        log.Fatalf("Nie udało się załadować konfiguracji: %v", err)
    }
    
    // Użyj konfiguracji
    fmt.Printf("Serwer: %s:%d\n", cfg.ServerHost, cfg.ServerPort)
    fmt.Printf("Tryb debugowania: %v\n", cfg.Debug)
}
```

### Format tagu struktury

Biblioteka używa tagu struktury `config` w następującym formacie:

```go
`config:"env=ENV_VAR_NAME,default=default_value,required=true"`
```

- `env`: Nazwa zmiennej środowiskowej, z której zostanie załadowana wartość
- `default`: Wartość domyślna, która zostanie użyta, jeśli zmienna środowiskowa nie jest ustawiona
- `required`: Ustawione na "true", aby oznaczyć pole jako wymagane (zwróci błąd, jeśli nie podano wartości)

Jeśli klucz `env` nie jest określony, nazwa pola w górnym rejestrze zostanie użyta jako nazwa zmiennej środowiskowej.

**Uwaga**: Jeśli pole jest oznaczone jako wymagane, ale ma wartość domyślną, wartość domyślna zostanie użyta, jeśli zmienna środowiskowa nie jest ustawiona, i nie zostanie zwrócony błąd.

### Obsługiwane typy

Biblioteka obsługuje następujące typy pól:

- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `bool`
- `time.Time` (format RFC3339, np. "2023-01-02T15:04:05Z")
- `time.Duration` (format czasu Go, np. "5s", "1h30m")
- `struct` (zagnieżdżone struktury)

### Zagnieżdżone struktury

Biblioteka obsługuje zagnieżdżone struktury dla lepszej organizacji konfiguracji. Możesz definiować zagnieżdżone struktury, aby grupować powiązane ustawienia konfiguracyjne:

```go
// ServerConfig reprezentuje konfigurację specyficzną dla serwera
type ServerConfig struct {
    Port int    `config:"env=SERVER_PORT,default=8080"`
    Host string `config:"env=SERVER_HOST,default=localhost"`
}

// DatabaseConfig reprezentuje konfigurację specyficzną dla bazy danych
type DatabaseConfig struct {
    Host     string `config:"env=DB_HOST,default=localhost"`
    Port     int    `config:"env=DB_PORT,default=5432"`
    User     string `config:"env=DB_USER,default=postgres"`
    Password string `config:"env=DB_PASSWORD,default=secret"`
    Name     string `config:"env=DB_NAME,default=mydb"`
}

// AppConfig reprezentuje konfigurację aplikacji z zagnieżdżonymi strukturami
type AppConfig struct {
    // Zagnieżdżone struktury konfiguracyjne
    Server   ServerConfig
    Database DatabaseConfig
    
    // Ustawienia aplikacji
    Debug bool `config:"env=DEBUG,default=false"`
}
```

Podczas korzystania z zagnieżdżonych struktur:

1. Biblioteka rekurencyjnie przetwarza wszystkie zagnieżdżone struktury
2. Każde pole w zagnieżdżonej strukturze może mieć własny tag config
3. Pola struktury bez tagów config nadal będą przetwarzane rekurencyjnie
4. Pozwala to na lepszą organizację ustawień konfiguracyjnych poprzez grupowanie powiązanych ustawień

### Uruchamianie z zmiennymi środowiskowymi

Możesz ustawić zmienne środowiskowe podczas uruchamiania aplikacji:

```bash
SERVER_PORT=9090 SERVER_HOST=0.0.0.0 DEBUG=true go run main.go
```

Lub ustawić je w swoim środowisku:

```bash
export SERVER_PORT=9090
export SERVER_HOST=0.0.0.0
export DEBUG=true
go run main.go
```

## Wymagane pola

Możesz oznaczyć pola jako wymagane, aby zapewnić, że mają wartości. Jeśli wymagane pole nie ma wartości ze zmiennej środowiskowej lub wartości domyślnej, zostanie zwrócony błąd.

```go
type Config struct {
    // To pole jest wymagane i musi być ustawione przez zmienną środowiskową
    APIKey string `config:"env=API_KEY,required=true"`
    
    // To pole jest wymagane, ale ma wartość domyślną, więc zawsze będzie miało wartość
    Timeout int `config:"env=TIMEOUT,default=30,required=true"`
    
    // To pole jest opcjonalne
    Debug bool `config:"env=DEBUG,default=false"`
}
```

Przykład obsługi błędów wymaganych pól:

```go
cfg := &Config{}
err := envconfig.Load(cfg)
if err != nil {
    // Sprawdź, czy to błąd wymaganego pola
    var reqErr *envconfig.RequiredFieldError
    if errors.As(err, &reqErr) {
        log.Fatalf("Brakujące wymagane pole: %s (zmienna środowiskowa: %s)", 
            reqErr.FieldName, reqErr.EnvName)
    }
    log.Fatalf("Nie udało się załadować konfiguracji: %v", err)
}
```

## Obsługa błędów

Biblioteka zapewnia szczegółowe raportowanie błędów walidacji i parsowania:

1. **RequiredFieldError**: Zwracany, gdy wymagane pole nie ma wartości
   - Zawiera nazwę pola i nazwę zmiennej środowiskowej

2. **ParseError**: Zwracany, gdy wartość nie może być sparsowana do docelowego typu
   - Zawiera nazwę pola, typ pola, wartość i podstawowy błąd

3. **ErrNotStruct**: Zwracany, gdy parametr konfiguracji nie jest wskaźnikiem do struktury

4. **ErrUnsupportedFieldType**: Zwracany, gdy pole ma nieobsługiwany typ

Przykład obsługi różnych typów błędów:

```go
cfg := &Config{}
err := envconfig.Load(cfg)
if err != nil {
    // Sprawdź konkretne typy błędów
    var reqErr *envconfig.RequiredFieldError
    var parseErr *envconfig.ParseError
    
    switch {
    case errors.As(err, &reqErr):
        log.Fatalf("Brakujące wymagane pole: %s (zmienna środowiskowa: %s)", 
            reqErr.FieldName, reqErr.EnvName)
    case errors.As(err, &parseErr):
        log.Fatalf("Nie udało się sparsować wartości '%s' jako %s dla pola '%s': %v", 
            parseErr.Value, parseErr.FieldType, parseErr.FieldName, parseErr.Err)
    case errors.Is(err, envconfig.ErrNotStruct):
        log.Fatalf("Konfiguracja musi być wskaźnikiem do struktury")
    case errors.Is(err, envconfig.ErrUnsupportedFieldType):
        log.Fatalf("Konfiguracja zawiera pole z nieobsługiwanym typem")
    default:
        log.Fatalf("Nie udało się załadować konfiguracji: %v", err)
    }
}
```

## Przykłady

Zobacz plik `main.go`, aby zobaczyć kompletny przykład użycia biblioteki.

## Changelog

### 2025-07-29
- **KLUCZOWA ZMIANA**: Zmieniono nazwę pakietu z `config` na `envconfig`. Ta zmiana wymaga aktualizacji importów w istniejącym kodzie.
- Nazwa tagu struktury pozostaje `config` dla zachowania kompatybilności wstecznej.

## Licencja

MIT