// Package tests содержит модульные тесты для проверки работы парсера конфигурации и графа.
package tests

import (
	"os"
	"path/filepath"
	"project2/src/config"
	"testing"
)

const (
	expectedMaxDepth = 3
	tempFilePerm     = 0600
)

// createTempConfig создает временный файл конфигурации во временной папке теста и возвращает путь к нему.
func createTempConfig(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test_config.toml")
	if err := os.WriteFile(path, []byte(content), tempFilePerm); err != nil {
		t.Fatalf("не удалось создать временный конфиг: %v", err)
	}
	return path
}

// TestValidConfig проверяет успешную загрузку и разбор корректно заполненного файла конфигурации.
func TestValidConfig(t *testing.T) {
	content := `
package_name = "tokio"
repo_url = "https://crates.io"
test_mode = false
output_file = "out.svg"
ascii_tree = true
max_depth = 3
`
	path := createTempConfig(t, content)
	cfg, err := config.LoadFromFile(path)
	if err != nil {
		t.Fatalf("ожидался успешный парсинг, получена ошибка: %v", err)
	}
	if cfg.PackageName != "tokio" || cfg.MaxDepth != expectedMaxDepth {
		t.Errorf("некорректные значения полей: %+v", cfg)
	}
}

// TestConfigValidationErrors проверяет корректную обработку ошибок при неверных или пропущенных параметрах.
func TestConfigValidationErrors(t *testing.T) {
	cases := []struct {
		name    string
		content string
	}{
		{
			name: "empty_package",
			content: "package_name = \"\"\nrepo_url = \"u\"\ntest_mode = false\n" +
				"output_file = \"o\"\nascii_tree = false\nmax_depth = 1",
		},
		{
			name: "invalid_bool",
			content: "package_name = \"p\"\nrepo_url = \"u\"\ntest_mode = yes\n" +
				"output_file = \"o\"\nascii_tree = false\nmax_depth = 1",
		},
		{
			name: "negative_depth",
			content: "package_name = \"p\"\nrepo_url = \"u\"\ntest_mode = false\n" +
				"output_file = \"o\"\nascii_tree = false\nmax_depth = 0",
		},
		{
			name:    "missing_field",
			content: "package_name = \"p\"\nrepo_url = \"u\"",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := createTempConfig(t, tc.content)
			if _, err := config.LoadFromFile(path); err == nil {
				t.Errorf("ожидалась ошибка валидации для кейса %s", tc.name)
			}
		})
	}
}
