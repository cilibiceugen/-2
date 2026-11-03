package tests

import (
	"project2/src/cargo"
	"slices"
	"testing"
)

// TestBuildIndexURL проверяет корректность построения URL для пакетов с разной длиной имени.
func TestBuildIndexURL(t *testing.T) {
	base := "https://crates.io"
	cases := []struct {
		pkg    string
		suffix string
	}{
		{"a", "/1/a"},
		{"io", "/2/io"},
		{"log", "/3/l/log"},
		{"serde", "/se/rd/serde"},
	}

	for _, tc := range cases {
		got, err := cargo.BuildIndexURL(base+"/", tc.pkg)
		if err != nil {
			t.Fatalf("неожиданная ошибка для %s: %v", tc.pkg, err)
		}
		expected := base + tc.suffix
		if got != expected {
			t.Errorf("для %s ожидалось %s, получено %s", tc.pkg, expected, got)
		}
	}
}

// TestParseLatestDependencies проверяет извлечение прямых зависимостей из последней версии JSON-индекса.
func TestParseLatestDependencies(t *testing.T) {
	sample := `{"name":"serde","vers":"1.0.0","deps":[{"name":"old_dep","kind":"normal"}]}
{"name":"serde","vers":"1.0.1","deps":[{"name":"serde_derive","kind":"normal"},` +
		`{"name":"syn","kind":"normal"},{"name":"trybuild","kind":"dev"}]}`

	deps, err := cargo.ParseLatestDependencies(sample)
	if err != nil {
		t.Fatalf("ошибка парсинга JSON: %v", err)
	}

	expected := []string{"serde_derive", "syn"}
	if !slices.Equal(deps, expected) {
		t.Errorf("ожидалось %v, получено %v", expected, deps)
	}
}
