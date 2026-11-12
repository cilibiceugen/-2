// Package visual отвечает за генерацию графа в формате Graphviz DOT, сохранение SVG и вывод ASCII-дерева.
package visual

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

const (
	svgFilePerm     = 0644
	svgNodeStepY    = 40
	svgPadding      = 30
	svgMinWidth     = 400
	svgCharWidth    = 9
	lastChildOffset = 1
)

// BuildDOT формирует текстовое описание направленного графа зависимостей на языке Graphviz.
func BuildDOT(root string, adj map[string][]string) string {
	var buf strings.Builder
	buf.WriteString("digraph Dependencies {\n")
	buf.WriteString("  rankdir=LR;\n")
	buf.WriteString("  node [shape=box, style=rounded];\n")

	keys := sortedKeys(adj)
	if len(keys) == 0 {
		buf.WriteString(fmt.Sprintf("  \"%s\";\n", root))
	}
	writeDOTEdges(&buf, keys, adj)

	buf.WriteString("}\n")
	return buf.String()
}

// writeDOTEdges записывает узлы и ребра графа в буфер в формате DOT.
func writeDOTEdges(buf *strings.Builder, keys []string, adj map[string][]string) {
	for _, parent := range keys {
		deps := adj[parent]
		if len(deps) == 0 {
			buf.WriteString(fmt.Sprintf("  \"%s\";\n", parent))
			continue
		}
		for _, child := range deps {
			buf.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", parent, child))
		}
	}
}

// SaveSVG сохраняет изображение графа в формате SVG с помощью Graphviz или встроенного генератора.
func SaveSVG(outputPath, dotContent string) error {
	svgBytes, err := renderWithGraphviz(dotContent)
	if err != nil {
		svgBytes = renderFallbackSVG(dotContent)
	}
	if err := os.WriteFile(outputPath, svgBytes, svgFilePerm); err != nil {
		return fmt.Errorf("ошибка записи SVG-файла %s: %w", outputPath, err)
	}
	return nil
}

// renderWithGraphviz вызывает внешнюю утилиту dot для преобразования DOT-текста в SVG.
func renderWithGraphviz(dotContent string) ([]byte, error) {
	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(dotContent)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// renderFallbackSVG создает базовый векторный SVG-файл с описанием ребер, если Graphviz не установлен.
func renderFallbackSVG(dotContent string) []byte {
	lines := strings.Split(strings.TrimSpace(dotContent), "\n")
	height := len(lines)*svgNodeStepY + svgPadding
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("<svg xmlns=\"http://w3.org\" width=\"%d\" height=\"%d\">\n",
		svgMinWidth*svgCharWidth/svgNodeStepY+svgMinWidth, height))
	for idx, line := range lines {
		yPos := (idx + lastChildOffset) * svgNodeStepY
		escaped := strings.ReplaceAll(line, "\"", "&quot;")
		buf.WriteString(fmt.Sprintf("  <text x=\"20\" y=\"%d\" font-family=\"monospace\">%s</text>\n", yPos, escaped))
	}
	buf.WriteString("</svg>\n")
	return []byte(buf.String())
}

// BuildASCIITree строит строковое представление зависимостей в виде псевдографического ASCII-дерева.
func BuildASCIITree(root string, adj map[string][]string) string {
	var buf strings.Builder
	buf.WriteString(root + "\n")
	visited := map[string]bool{root: true}
	renderTreeBranch(&buf, root, "", adj, visited)
	return buf.String()
}

// renderTreeBranch рекурсивно форматирует дочерние узлы ASCII-дерева с защитой от циклов.
func renderTreeBranch(
	buf *strings.Builder,
	node string,
	prefix string,
	adj map[string][]string,
	visited map[string]bool,
) {
	children := adj[node]
	for i, child := range children {
		isLast := i == len(children)-lastChildOffset
		connector, nextPrefix := branchSymbols(isLast)
		if visited[child] {
			buf.WriteString(fmt.Sprintf("%s%s%s (*)\n", prefix, connector, child))
			continue
		}
		buf.WriteString(fmt.Sprintf("%s%s%s\n", prefix, connector, child))
		visited[child] = true
		renderTreeBranch(buf, child, prefix+nextPrefix, adj, visited)
		delete(visited, child)
	}
}

// branchSymbols возвращает символы ветвления для текущего и вложенного уровней дерева.
func branchSymbols(isLast bool) (string, string) {
	if isLast {
		return "└── ", "    "
	}
	return "├── ", "│   "
}

// sortedKeys возвращает отсортированный список ключей карты смежности.
func sortedKeys(adj map[string][]string) []string {
	keys := make([]string, 0, len(adj))
	for k := range adj {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
