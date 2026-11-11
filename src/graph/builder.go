// Package graph реализует построение и анализ транзитивного графа зависимостей.
package graph

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	initialBuildDepth  = 0
	testLinePartsCount = 2
	testLineSeparator  = ":"
	testDepsSeparator  = ","
)

// DepFetcher определяет функцию получения прямых зависимостей для указанного пакета.
type DepFetcher func(repoURL, pkgName string) ([]string, error)

// BuildFromCargo строит транзитивный граф зависимостей из индекса Cargo с помощью рекурсивного BFS.
func BuildFromCargo(
	repoURL string,
	rootPkg string,
	maxDepth int,
	fetch DepFetcher,
) (map[string][]string, error) {
	adj := make(map[string][]string)
	visited := map[string]bool{rootPkg: true}

	fetcher := func(pkg string) ([]string, error) {
		return fetch(repoURL, pkg)
	}

	if err := bfsBuildRecursive([]string{rootPkg}, initialBuildDepth, maxDepth, visited, adj, fetcher); err != nil {
		return nil, err
	}
	return adj, nil
}

// BuildFromTestFile читает описание тестового графа из файла и строит граф зависимостей рекурсивным BFS.
func BuildFromTestFile(filePath, rootPkg string, maxDepth int) (map[string][]string, error) {
	rawGraph, err := parseTestRepository(filePath)
	if err != nil {
		return nil, err
	}
	if _, exists := rawGraph[rootPkg]; !exists {
		return nil, fmt.Errorf("корневой пакет %s не найден в тестовом файле", rootPkg)
	}

	adj := make(map[string][]string)
	visited := map[string]bool{rootPkg: true}
	fetcher := func(pkg string) ([]string, error) {
		return rawGraph[pkg], nil
	}

	if err := bfsBuildRecursive([]string{rootPkg}, initialBuildDepth, maxDepth, visited, adj, fetcher); err != nil {
		return nil, err
	}
	return adj, nil
}

// bfsBuildRecursive выполняет рекурсивный обход графа в ширину с учетом maxDepth и защитой от циклов.
func bfsBuildRecursive(
	level []string,
	depth int,
	maxDepth int,
	visited map[string]bool,
	adj map[string][]string,
	fetch func(string) ([]string, error),
) error {
	if len(level) == 0 || depth >= maxDepth {
		return nil
	}

	var nextLevel []string
	for _, pkg := range level {
		deps, err := fetch(pkg)
		if err != nil {
			return err
		}
		adj[pkg] = appendSortedUnique(adj[pkg], deps)
		nextLevel = collectUnvisited(deps, visited, nextLevel)
	}
	return bfsBuildRecursive(nextLevel, depth+depthIncrement, maxDepth, visited, adj, fetch)
}

// collectUnvisited отмечает непосещенные узлы и добавляет их в очередь следующего уровня обхода.
func collectUnvisited(deps []string, visited map[string]bool, nextLevel []string) []string {
	for _, dep := range deps {
		if !visited[dep] {
			visited[dep] = true
			nextLevel = append(nextLevel, dep)
		}
	}
	return nextLevel
}

// appendSortedUnique объединяет списки зависимостей, сортирует их и удаляет дубликаты.
func appendSortedUnique(existing, added []string) []string {
	combined := append(existing, added...)
	slices.Sort(combined)
	return slices.Compact(combined)
}

// parseTestRepository считывает тестовый файл формата "A: B, C" в карту смежности.
func parseTestRepository(path string) (map[string][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия тестового репозитория: %w", err)
	}
	defer file.Close()

	result := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := parseTestLine(scanner.Text(), result); err != nil {
			return nil, err
		}
	}
	return result, scanner.Err()
}

// parseTestLine разбирает одну строку вида "A: B, C" и добавляет узел в карту смежности.
func parseTestLine(rawLine string, result map[string][]string) error {
	line := strings.TrimSpace(rawLine)
	if line == "" {
		return nil
	}
	parts := strings.SplitN(line, testLineSeparator, testLinePartsCount)
	if len(parts) != testLinePartsCount {
		return fmt.Errorf("некорректный формат строки в тестовом файле: %s", line)
	}
	node := strings.TrimSpace(parts[0])
	result[node] = splitDependencies(parts[1])
	return nil
}

// splitDependencies разбивает правую часть строки тестового файла по запятым на имена пакетов.
func splitDependencies(rawDeps string) []string {
	trimmed := strings.TrimSpace(rawDeps)
	if trimmed == "" {
		return nil
	}
	var deps []string
	for _, item := range strings.Split(trimmed, testDepsSeparator) {
		if name := strings.TrimSpace(item); name != "" {
			deps = append(deps, name)
		}
	}
	return deps
}
