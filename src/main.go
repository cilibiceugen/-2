// Package main запускает CLI-приложение для анализа и визуализации графа зависимостей.
package main

import (
	"flag"
	"fmt"
	"os"
	"project2/src/cargo"
	"project2/src/config"
)

const exitErrorCode = 1

// main разбирает флаг командной строки, загружает конфигурацию и выводит прямые зависимости пакета.
func main() {
	configPath := flag.String("config", "config.toml", "Путь к файлу конфигурации TOML")
	flag.Parse()

	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка конфигурации: %v\n", err)
		os.Exit(exitErrorCode)
	}

	deps, err := cargo.FetchDirectDependencies(cfg.RepoURL, cfg.PackageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка получения зависимостей: %v\n", err)
		os.Exit(exitErrorCode)
	}

	printDependencies(cfg.PackageName, deps)
}

// printDependencies выводит на экран список прямых зависимостей указанного пакета.
func printDependencies(pkgName string, deps []string) {
	fmt.Printf("Прямые зависимости пакета %s (%d):\n", pkgName, len(deps))
	if len(deps) == 0 {
		fmt.Println("  (зависимости отсутствуют)")
		return
	}
	for _, dep := range deps {
		fmt.Printf("  - %s\n", dep)
	}
}
