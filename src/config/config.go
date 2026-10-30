// Package config отвечает за чтение файла настроек TOML и проверку всех параметров.
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	expectedPartsCount = 2
	minValidDepth      = 1
	commentPrefix      = "#"
	keyValueSeparator  = "="
)

// AppConfig хранит все шесть параметров, которые пользователь задает в файле конфигурации.
type AppConfig struct {
	PackageName string
	RepoURL     string
	TestMode    bool
	OutputFile  string
	ASCIITree   bool
	MaxDepth    int
}

// LoadFromFile открывает файл конфигурации, читает строки в словарь и проверяет все значения.
func LoadFromFile(path string) (AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return AppConfig{}, fmt.Errorf("ошибка открытия файла конфигурации: %w", err)
	}
	defer file.Close()

	rawMap, err := parseTOMLLines(file)
	if err != nil {
		return AppConfig{}, err
	}

	return buildAndValidate(rawMap)
}

// parseTOMLLines построчно читает открытый файл, пропуская пустые строки и комментарии с символом решетки.
func parseTOMLLines(file *os.File) (map[string]string, error) {
	result := make(map[string]string)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, commentPrefix) {
			continue
		}
		key, val, err := splitKeyValue(line, lineNum)
		if err != nil {
			return nil, err
		}
		result[key] = val
	}
	return result, scanner.Err()
}

// splitKeyValue делит строку по первому знаку равенства на очищенные от пробелов ключ и значение.
func splitKeyValue(line string, lineNum int) (string, string, error) {
	parts := strings.SplitN(line, keyValueSeparator, expectedPartsCount)
	if len(parts) != expectedPartsCount {
		return "", "", fmt.Errorf("синтаксическая ошибка в строке %d: %s", lineNum, line)
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	return key, val, nil
}

// buildAndValidate проверяет наличие всех обязательных полей и переводит их в нужные типы данных.
func buildAndValidate(raw map[string]string) (AppConfig, error) {
	var cfg AppConfig
	var err error

	if cfg.PackageName, err = parseStringParam(raw, "package_name"); err != nil {
		return cfg, err
	}
	if cfg.RepoURL, err = parseStringParam(raw, "repo_url"); err != nil {
		return cfg, err
	}
	if cfg.OutputFile, err = parseStringParam(raw, "output_file"); err != nil {
		return cfg, err
	}
	if cfg.TestMode, err = parseBoolParam(raw, "test_mode"); err != nil {
		return cfg, err
	}
	if cfg.ASCIITree, err = parseBoolParam(raw, "ascii_tree"); err != nil {
		return cfg, err
	}
	if cfg.MaxDepth, err = parseDepthParam(raw, "max_depth"); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// parseStringParam проверяет, что значение существует, заключено в двойные кавычки и не пустое.
func parseStringParam(raw map[string]string, key string) (string, error) {
	val, exists := raw[key]
	if !exists {
		return "", fmt.Errorf("отсутствует обязательный параметр: %s", key)
	}
	if !strings.HasPrefix(val, "\"") || !strings.HasSuffix(val, "\"") {
		return "", fmt.Errorf("параметр %s должен быть строкой в двойных кавычках", key)
	}
	cleaned := strings.Trim(val, "\"")
	if cleaned == "" {
		return "", fmt.Errorf("параметр %s не может быть пустым", key)
	}
	return cleaned, nil
}

// parseBoolParam проверяет наличие параметра и переводит строку true или false в логический тип.
func parseBoolParam(raw map[string]string, key string) (bool, error) {
	val, exists := raw[key]
	if !exists {
		return false, fmt.Errorf("отсутствует обязательный параметр: %s", key)
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return false, fmt.Errorf("параметр %s должен быть true или false", key)
	}
	return parsed, nil
}

// parseDepthParam проверяет, что глубина анализа является целым числом не меньше единицы.
func parseDepthParam(raw map[string]string, key string) (int, error) {
	val, exists := raw[key]
	if !exists {
		return 0, fmt.Errorf("отсутствует обязательный параметр: %s", key)
	}
	depth, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("параметр %s должен быть целым числом", key)
	}
	if depth < minValidDepth {
		return 0, errors.New("параметр max_depth должен быть больше нуля")
	}
	return depth, nil
}
