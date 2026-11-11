package tests

import (
	"project2/src/graph"
	"slices"
	"testing"
)

const (
	testFullDepth    = 5
	testLimitedDepth = 1
)

// TestFindReverseDependencies проверяет поиск прямых и транзитивных обратных зависимостей с циклами.
func TestFindReverseDependencies(t *testing.T) {
	adj := map[string][]string{
		"A": {"B", "C"},
		"B": {"C", "D"},
		"C": {"A", "E"},
		"D": {"F"},
		"E": {},
		"F": {},
	}

	got := graph.FindReverseDependencies(adj, "F", testFullDepth)
	expected := []string{"A", "B", "C", "D"}
	if !slices.Equal(got, expected) {
		t.Errorf("для F ожидалось %v, получено %v", expected, got)
	}
}

// TestReverseDepthLimit проверяет ограничение глубины при поиске обратных зависимостей.
func TestReverseDepthLimit(t *testing.T) {
	adj := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {},
	}

	got := graph.FindReverseDependencies(adj, "C", testLimitedDepth)
	expected := []string{"B"}
	if !slices.Equal(got, expected) {
		t.Errorf("на глубине 1 ожидалось %v, получено %v", expected, got)
	}
}
