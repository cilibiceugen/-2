package tests

import (
	"os"
	"path/filepath"
	"project2/src/graph"
	"slices"
	"testing"
)

const (
	graphFullDepth    = 4
	graphLimitedDepth = 1
	graphFilePerm     = 0600
)

// createTempRepoFile создает временный текстовый файл тестового репозитория и возвращает путь к нему.
func createTempRepoFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "repo.txt")
	if err := os.WriteFile(path, []byte(content), graphFilePerm); err != nil {
		t.Fatalf("ошибка создания временного файла репозитория: %v", err)
	}
	return path
}

// TestBuildFromTestFileWithCycles проверяет построение транзитивного графа и защиту от циклов.
func TestBuildFromTestFileWithCycles(t *testing.T) {
	content := "A: B, C\nB: C, D\nC: A, E\nD:\nE:\n"
	path := createTempRepoFile(t, content)

	adj, err := graph.BuildFromTestFile(path, "A", graphFullDepth)
	if err != nil {
		t.Fatalf("ошибка построения графа: %v", err)
	}

	if !slices.Equal(adj["A"], []string{"B", "C"}) {
		t.Errorf("неверные зависимости для A: %v", adj["A"])
	}
	if !slices.Equal(adj["C"], []string{"A", "E"}) {
		t.Errorf("неверные зависимости для C: %v", adj["C"])
	}
}

// TestBuildDepthLimit проверяет ограничение обхода графа параметром максимальной глубины.
func TestBuildDepthLimit(t *testing.T) {
	content := "A: B\nB: C\nC: D\nD:\n"
	path := createTempRepoFile(t, content)

	adj, err := graph.BuildFromTestFile(path, "A", graphLimitedDepth)
	if err != nil {
		t.Fatalf("ошибка построения графа: %v", err)
	}

	if _, exists := adj["B"]; exists {
		t.Errorf("узел B не должен раскрываться при глубине 1: %v", adj)
	}
}

// TestBuildFromCargoMock проверяет построение графа через пользовательскую функцию получения зависимостей.
func TestBuildFromCargoMock(t *testing.T) {
	mockData := map[string][]string{
		"root": {"mid"},
		"mid":  {"leaf"},
		"leaf": {},
	}
	fetcher := func(_, pkg string) ([]string, error) {
		return mockData[pkg], nil
	}

	adj, err := graph.BuildFromCargo("https://crates.io", "root", graphFullDepth, fetcher)
	if err != nil {
		t.Fatalf("ошибка построения графа: %v", err)
	}
	if !slices.Equal(adj["mid"], []string{"leaf"}) {
		t.Errorf("ожидалась транзитивная зависимость leaf, получено: %v", adj["mid"])
	}
}
