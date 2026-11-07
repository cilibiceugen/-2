// Package graph реализует построение графа зависимостей рекурсивным алгоритмом BFS.
package graph

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	initialDepth      = 0
	depthStep         = 1
	expectedPairParts = 2
	pairSeparator     = ":"
	depsSeparator     = ","
)

// DependencyProvider определяет источник получения прямых зависимостей для указанного пакета.
type DependencyProvider func(pkg string) ([]string, error)

// DependencyGraph представляет построенный граф зависимостей в виде списков смежности.
type DependencyGraph struct {
	Root  string
	Edges map[string][]string
}

// BuildGraph строит граф зависимостей от корневого пакета с помощью рекурсивного BFS.
func BuildGraph(root string, maxDepth int, provider DependencyProvider) (DependencyGraph, error) {
	graph := DependencyGraph{
		Root:  root,
		Edges: make(map[string][]string),
	}
	visited := map[string]bool{root: true}
	currentLevel := []string{root}

	if err := recursiveBFS(&graph, currentLevel, visited, initialDepth, maxDepth, provider); err != nil {
		return DependencyGraph{}, err
	}
	return graph, nil
}

// recursiveBFS выполняет рекурсивный поиск в ширину по уровням глубины от текущего слоя узлов.
func recursiveBFS(
	graph *DependencyGraph,
	level []string,
	visited map[string]bool,
	depth int,
	maxDepth int,
	provider DependencyProvider,
) error {
	if len(level) == 0 || depth >= maxDepth {
		return nil
	}

	nextLevel, err := expandLevel(graph, level, visited, provider)
	if err != nil {
		return err
	}

	return recursiveBFS(graph, nextLevel, visited, depth+depthStep, maxDepth, provider)
}

// expandLevel запрашивает зависимости для всех пакетов текущего уровня и формирует следующий слой BFS.
func expandLevel(
	graph *DependencyGraph,
	level []string,
	visited map[string]bool,
	provider DependencyProvider,
) ([]string, error) {
	var nextLevel []string

	for _, pkg := range level {
		deps, err := provider(pkg)
		if err != nil {
			return nil, err
		}
		graph.Edges[pkg] = deps
		nextLevel = appendUnvisited(nextLevel, deps, visited)
	}
	return nextLevel, nil
}

// appendUnvisited добавляет еще не посещенные пакеты в следующий слой обхода и помечает их.
func appendUnvisited(nextLevel []string, deps []string, visited map[string]bool) []string {
	for _, dep := range deps {
		if !visited[dep] {
			visited[dep] = true
			nextLevel = append(nextLevel, dep)
		}
	}
	return nextLevel
}

// LoadTestRepository читает тестовый граф из файла формата "A: B, C" с именами из заглавных букв.
func LoadTestRepository(filePath string) (DependencyProvider, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия тестового репозитория: %w", err)
	}
	defer file.Close()

	repoMap, err := parseTestLines(file)
	if err != nil {
		return nil, err
	}

	return func(pkg string) ([]string, error) {
		return repoMap[pkg], nil
	}, nil
}

// parseTestLines построчно разбирает описание зависимостей тестового репозитория в словарь.
func parseTestLines(file *os.File) (map[string][]string, error) {
	repo := make(map[string][]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		pkg, deps, err := parseTestLine(line)
		if err != nil {
			return nil, err
		}
		repo[pkg] = deps
	}
	return repo, scanner.Err()
}

// parseTestLine разбирает одну строку вида "A: B, C" и проверяет корректность имен пакетов.
func parseTestLine(line string) (string, []string, error) {
	parts := strings.SplitN(line, pairSeparator, expectedPairParts)
	if len(parts) != expectedPairParts {
		return "", nil, fmt.Errorf("неверный формат строки тестового графа: %s", line)
	}

	pkg := strings.TrimSpace(parts[0])
	if !isUpperLatin(pkg) {
		return "", nil, fmt.Errorf("имя пакета в тестовом режиме должно быть из заглавных букв: %s", pkg)
	}

	deps, err := parseTargetList(parts[1])
	if err != nil {
		return "", nil, err
	}
	return pkg, deps, nil
}

// parseTargetList разбивает правую часть строки по запятым и проверяет имена зависимостей.
func parseTargetList(rawTargets string) ([]string, error) {
	trimmed := strings.TrimSpace(rawTargets)
	if trimmed == "" {
		return nil, nil
	}

	items := strings.Split(trimmed, depsSeparator)
	var deps []string
	for _, item := range items {
		name := strings.TrimSpace(item)
		if !isUpperLatin(name) {
			return nil, fmt.Errorf("некорректное имя зависимости в тесте: %s", name)
		}
		deps = append(deps, name)
	}
	return deps, nil
}

// isUpperLatin проверяет, что строка не пуста и состоит только из заглавных латинских букв A-Z.
func isUpperLatin(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	return true
}
