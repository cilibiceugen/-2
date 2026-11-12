// Package main запускает CLI-приложение для анализа и визуализации графа зависимостей.
package main

import (
	"flag"
	"fmt"
	"os"
	"project2/src/cargo"
	"project2/src/config"
	"project2/src/graph"
	"project2/src/visual"
	"slices"
	"strings"
)

const exitErrorCode = 1

// main разбирает флаги, строит граф зависимостей, выводит дерево и сохраняет SVG-визуализацию.
func main() {
	configPath := flag.String("config", "config.toml", "Путь к файлу конфигурации TOML")
	reverseTarget := flag.String("reverse", "", "Имя пакета для поиска обратных зависимостей")
	flag.Parse()

	cfg, err := loadAndBuild(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(exitErrorCode)
	}

	runVisualization(cfg, *reverseTarget)
}

// loadAndBuild загружает конфигурацию из файла.
func loadAndBuild(configPath *string) (config.AppConfig, error) {
	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		return cfg, fmt.Errorf("ошибка конфигурации: %w", err)
	}
	return cfg, nil
}

// runVisualization строит граф, выводит результаты в консоль и сохраняет SVG-файл.
func runVisualization(cfg config.AppConfig, reverseFlag string) {
	adj, err := loadDependencyGraph(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка построения графа: %v\n", err)
		os.Exit(exitErrorCode)
	}

	printOutputs(cfg, adj, reverseFlag)
	dotText := visual.BuildDOT(cfg.PackageName, adj)
	if err := visual.SaveSVG(cfg.OutputFile, dotText); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(exitErrorCode)
	}
	fmt.Printf("\nГраф успешно сохранен в файл: %s\n", cfg.OutputFile)
}

// printOutputs выводит граф, обратные зависимости и ASCII-дерево при включенном флаге.
func printOutputs(cfg config.AppConfig, adj map[string][]string, reverseFlag string) {
	printGraph(cfg.PackageName, adj)
	target := resolveReverseTarget(reverseFlag, cfg.PackageName)
	revDeps := graph.FindReverseDependencies(adj, target, cfg.MaxDepth)
	printReverseDeps(target, revDeps)

	if cfg.ASCIITree {
		fmt.Printf("\nASCII-дерево зависимостей:\n%s", visual.BuildASCIITree(cfg.PackageName, adj))
	}
}

// loadDependencyGraph строит граф из тестового файла или из индекса Cargo в зависимости от test_mode.
func loadDependencyGraph(cfg config.AppConfig) (map[string][]string, error) {
	if cfg.TestMode {
		return graph.BuildFromTestFile(cfg.RepoURL, cfg.PackageName, cfg.MaxDepth)
	}
	return graph.BuildFromCargo(cfg.RepoURL, cfg.PackageName, cfg.MaxDepth, cargo.FetchDirectDependencies)
}

// resolveReverseTarget возвращает имя пакета из флага или корневой пакет конфигурации по умолчанию.
func resolveReverseTarget(flagVal, defaultPkg string) string {
	if strings.TrimSpace(flagVal) != "" {
		return strings.TrimSpace(flagVal)
	}
	return defaultPkg
}

// printGraph выводит список смежности построенного графа зависимостей в отсортированном виде.
func printGraph(root string, adj map[string][]string) {
	fmt.Printf("Граф зависимостей для пакета %s:\n", root)
	keys := make([]string, 0, len(adj))
	for k := range adj {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("  %s -> %s\n", k, strings.Join(adj[k], ", "))
	}
}

// printReverseDeps выводит все пакеты, которые прямо или транзитивно зависят от target.
func printReverseDeps(target string, revDeps []string) {
	fmt.Printf("\nОбратные зависимости для пакета %s (%d):\n", target, len(revDeps))
	if len(revDeps) == 0 {
		fmt.Println("  (обратные зависимости отсутствуют)")
		return
	}
	for _, dep := range revDeps {
		fmt.Printf("  - %s\n", dep)
	}
}
