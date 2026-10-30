// Package main запускает CLI-приложение для анализа и визуализации графа зависимостей.
package main

import (
	"Project(2)/config"
	"flag"
	"fmt"
	"os"
)

const exitErrorCode = 1

// main разбирает флаг командной строки, загружает конфигурацию и выводит параметры на экран.
func main() {
	configPath := flag.String("config", "config.toml", "Путь к файлу конфигурации TOML")
	flag.Parse()

	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка конфигурации: %v\n", err)
		os.Exit(exitErrorCode)
	}

	printConfig(cfg)
}

// printConfig выводит все шесть прочитанных параметров конфигурации в формате ключ-значение.
func printConfig(cfg config.AppConfig) {
	fmt.Printf("package_name: %s\n", cfg.PackageName)
	fmt.Printf("repo_url: %s\n", cfg.RepoURL)
	fmt.Printf("test_mode: %t\n", cfg.TestMode)
	fmt.Printf("output_file: %s\n", cfg.OutputFile)
	fmt.Printf("ascii_tree: %t\n", cfg.ASCIITree)
	fmt.Printf("max_depth: %d\n", cfg.MaxDepth)
}
