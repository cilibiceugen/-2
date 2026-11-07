// Package main запускает CLI-приложение для анализа и визуализации графа зависимостей.
package main

import (
	"flag"
	"fmt"
	"os"
	"project2/src/cargo"
	"project2/src/config"
	"project2/src/graph"
	"sort"
	"strings"
)

const exitErrorCode = 1

// main разбирает флаг командной строки, загружает конфигурацию и строит граф зависимостей.
func main() {
	configPath := flag.String("config", "config.toml", "Путь к файлу конфигурации TOML")
	flag.Parse()

	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка конфигурации: %v\n", err)
		os.Exit(exitErrorCode)
	}

	provider, err := createProvider(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации источника данных: %v\n", err)
		os.Exit(exitErrorCode)
	}

	depGraph, err := graph.BuildGraph(cfg.PackageName, cfg.MaxDepth, provider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка построения графа: %v\n", err)
		os.Exit(exitErrorCode)
	}

	printGraph(depGraph)
}

// createProvider выбирает источник зависимостей в зависимости от флага тестового режима.
func createProvider(cfg config.AppConfig) (graph.DependencyProvider, error) {
	if cfg.TestMode {
		return graph.LoadTestRepository(cfg.RepoURL)
	}
	return func(pkg string) ([]string, error) {
		return cargo.FetchDirectDependencies(cfg.RepoURL, pkg)
	}, nil
}

// printGraph выводит построенный граф зависимостей в отсортированном виде.
func printGraph(depGraph graph.DependencyGraph) {
	fmt.Printf("Граф зависимостей для пакета %s:\n", depGraph.Root)
	nodes := sortedKeys(depGraph.Edges)
	for _, node := range nodes {
		deps := depGraph.Edges[node]
		if len(deps) == 0 {
			fmt.Printf("  %s -> (нет зависимостей)\n", node)
			continue
		}
		fmt.Printf("  %s -> %s\n", node, strings.Join(deps, ", "))
	}
}

// sortedKeys возвращает отсортированный список ключей словаря графа для стабильного вывода.
func sortedKeys(edges map[string][]string) []string {
	keys := make([]string, 0, len(edges))
	for k := range edges {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
