// Package cargo реализует загрузку и разбор метаданных пакетов Rust из индекса Cargo.
package cargo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	lenOneChar    = 1
	lenTwoChars   = 2
	lenThreeChars = 3
	firstCharIdx  = 1
	secondPairIdx = 4
	normalDepKind = "normal"
)

type packageVersion struct {
	Deps []rawDependency `json:"deps"`
}

type rawDependency struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// FetchDirectDependencies загружает метаданные пакета по URL индекса и возвращает список прямых зависимостей.
func FetchDirectDependencies(repoURL, pkgName string) ([]string, error) {
	fullURL, err := BuildIndexURL(repoURL, pkgName)
	if err != nil {
		return nil, err
	}

	body, err := downloadIndexFile(fullURL)
	if err != nil {
		return nil, err
	}

	return ParseLatestDependencies(body)
}

// BuildIndexURL формирует полный URL к файлу пакета в разреженном индексе Cargo по правилам длины имени.
func BuildIndexURL(baseURL, pkgName string) (string, error) {
	cleaned := strings.ToLower(strings.TrimSpace(pkgName))
	if cleaned == "" {
		return "", fmt.Errorf("имя пакета не может быть пустым")
	}

	base := strings.TrimRight(baseURL, "/")
	path := buildRelativePath(cleaned)
	return fmt.Sprintf("%s/%s", base, path), nil
}

// buildRelativePath вычисляет относительный путь в индексе Cargo в зависимости от количества символов в имени.
func buildRelativePath(name string) string {
	switch len(name) {
	case lenOneChar:
		return fmt.Sprintf("1/%s", name)
	case lenTwoChars:
		return fmt.Sprintf("2/%s", name)
	case lenThreeChars:
		return fmt.Sprintf("3/%s/%s", name[:firstCharIdx], name)
	default:
		return fmt.Sprintf("%s/%s/%s", name[:lenTwoChars], name[lenTwoChars:secondPairIdx], name)
	}
}

// downloadIndexFile выполняет HTTP GET запрос с заголовком User-Agent и возвращает тело ответа.
func downloadIndexFile(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса к %s: %w", url, err)
	}
	req.Header.Set("User-Agent", "project2-depgraph-student-tool")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка сетевого запроса к %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("репозиторий вернул статус %d для %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа сервера: %w", err)
	}
	return string(data), nil
}

// ParseLatestDependencies извлекает уникальные имена прямых зависимостей из последней версии в индексе Cargo.
func ParseLatestDependencies(content string) ([]string, error) {
	lastLine, err := extractLastNonEmptyLine(content)
	if err != nil {
		return nil, err
	}

	var ver packageVersion
	if err := json.Unmarshal([]byte(lastLine), &ver); err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON метаданных пакета: %w", err)
	}

	return filterDirectDeps(ver.Deps), nil
}

// extractLastNonEmptyLine находит последнюю непустую строку, соответствующую актуальной версии пакета.
func extractLastNonEmptyLine(content string) (string, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "" {
			return trimmed, nil
		}
	}
	return "", fmt.Errorf("файл метаданных пакета пуст")
}

// filterDirectDeps отбирает обычные зависимости (normal) и исключает повторяющиеся имена.
func filterDirectDeps(deps []rawDependency) []string {
	seen := make(map[string]bool)
	var result []string

	for _, dep := range deps {
		if dep.Kind != "" && dep.Kind != normalDepKind {
			continue
		}
		if !seen[dep.Name] {
			seen[dep.Name] = true
			result = append(result, dep.Name)
		}
	}
	return result
}
