package graph

import "slices"

const (
	initialReverseDepth = 0
	depthIncrement      = 1
)

// InvertGraph разворачивает все направленные ребра графа для поиска обратных зависимостей.
func InvertGraph(adj map[string][]string) map[string][]string {
	inverted := make(map[string][]string)
	for parent, children := range adj {
		if _, exists := inverted[parent]; !exists {
			inverted[parent] = nil
		}
		for _, child := range children {
			inverted[child] = append(inverted[child], parent)
		}
	}
	for node := range inverted {
		slices.Sort(inverted[node])
		inverted[node] = slices.Compact(inverted[node])
	}
	return inverted
}

// FindReverseDependencies находит все пакеты, прямо или косвенно зависящие от target, с помощью рекурсивного BFS.
func FindReverseDependencies(adj map[string][]string, target string, maxDepth int) []string {
	inverted := InvertGraph(adj)
	if _, exists := inverted[target]; !exists || maxDepth < depthIncrement {
		return nil
	}

	visited := map[string]bool{target: true}
	var result []string
	bfsReverseRecursive(inverted, []string{target}, initialReverseDepth, maxDepth, visited, &result)
	slices.Sort(result)
	return result
}

// bfsReverseRecursive выполняет рекурсивный обход в ширину по инвертированному графу с учетом глубины.
func bfsReverseRecursive(
	inverted map[string][]string,
	level []string,
	depth int,
	maxDepth int,
	visited map[string]bool,
	result *[]string,
) {
	if len(level) == 0 || depth >= maxDepth {
		return
	}
	var nextLevel []string
	for _, current := range level {
		for _, dependent := range inverted[current] {
			if !visited[dependent] {
				visited[dependent] = true
				*result = append(*result, dependent)
				nextLevel = append(nextLevel, dependent)
			}
		}
	}
	bfsReverseRecursive(inverted, nextLevel, depth+depthIncrement, maxDepth, visited, result)
}
