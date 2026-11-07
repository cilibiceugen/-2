package tests

import (
	"os"
	"path/filepath"
	"project2/src/graph"
	"slices"
	"testing"
)

const (
	depthOne   = 1
	depthTwo   = 2
	depthLarge = 10
)

// TestBuildGraphCyclicAndDepth проверяет рекурсивный BFS на графе с циклами и ограничением глубины.
func TestBuildGraphCyclicAndDepth(t *testing.T) {
	mockRepo := map[string][]string{
		"A": {"B", "C"},
		"B": {"D", "A"},
		"C": {"E"},
		"D": {"F"},
		"E": {},
		"F": {},
	}
	provider := func(pkg string) ([]string, error) {
		return mockRepo[pkg], nil
	}

	gDepth1, err := graph.BuildGraph("A", depthOne, provider)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if _, expandedB := gDepth1.Edges["B"]; expandedB {
		t.Errorf("при глубине 1 узел B не должен раскрываться")
	}

	gFull, err := graph.BuildGraph("A", depthLarge, provider)
	if err != nil {
		t.Fatalf("ошибка на циклическом графе: %v", err)
	}
	if !slices.Equal(gFull.Edges["D"], []string{"F"}) {
		t.Errorf("ожидалась зависимость D -> F, получено %v", gFull.Edges["D"])
	}
}

// TestLoadTestRepository проверяет чтение графа из текстового файла с заглавными буквами.
func TestLoadTestRepository(t *testing.T) {
	content := "A: B, C\nB: A\nC:\n"
	path := filepath.Join(t.TempDir(), "repo.txt")
	if err := os.WriteFile(path, []byte(content), tempFilePerm); err != nil {
		t.Fatalf("ошибка создания тестового файла: %v", err)
	}

	provider, err := graph.LoadTestRepository(path)
	if err != nil {
		t.Fatalf("ошибка загрузки тестового репозитория: %v", err)
	}

	g, err := graph.BuildGraph("A", depthTwo, provider)
	if err != nil {
		t.Fatalf("ошибка построения графа: %v", err)
	}
	if !slices.Equal(g.Edges["B"], []string{"A"}) {
		t.Errorf("некорректные ребра для B: %v", g.Edges["B"])
	}
}
