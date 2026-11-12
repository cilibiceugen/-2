package tests

import (
	"os"
	"path/filepath"
	"project2/src/visual"
	"strings"
	"testing"
)

// TestBuildDOTAndSVG проверяет генерацию кода Graphviz DOT и создание непустого SVG-файла.
func TestBuildDOTAndSVG(t *testing.T) {
	adj := map[string][]string{
		"A": {"B", "C"},
		"B": {},
		"C": {},
	}
	dot := visual.BuildDOT("A", adj)
	if !strings.Contains(dot, `"A" -> "B";`) || !strings.Contains(dot, `"A" -> "C";`) {
		t.Errorf("некорректный вывод DOT:\n%s", dot)
	}

	outPath := filepath.Join(t.TempDir(), "test_graph.svg")
	if err := visual.SaveSVG(outPath, dot); err != nil {
		t.Fatalf("ошибка сохранения SVG: %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil || !strings.Contains(string(data), "<svg") {
		t.Errorf("сгенерированный файл не является валидным SVG")
	}
}

// TestBuildASCIITree проверяет построение ASCII-дерева и пометку циклических зависимостей.
func TestBuildASCIITree(t *testing.T) {
	adj := map[string][]string{
		"A": {"B"},
		"B": {"A"},
	}
	tree := visual.BuildASCIITree("A", adj)
	if !strings.Contains(tree, "└── B") || !strings.Contains(tree, "└── A (*)") {
		t.Errorf("некорректное ASCII-дерево:\n%s", tree)
	}
}
