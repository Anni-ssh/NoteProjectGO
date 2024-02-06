package helperPath

import (
	"TestProject/internal/config"
	"fmt"
	"os"
	"path/filepath"
)

// RootDir возвращает абсолютный путь до корневой директории.
func RootDir() (string, error) {
	const operation = "helperFunc.RootDir"
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%s:%w", operation, err)
	}
	rootPath, err := filepath.Abs(currentDir)
	if err != nil {
		return "", fmt.Errorf("%s:%w", operation, err)
	}
	path := filepath.Dir(rootPath)
	fixPath := filepath.Dir(path)

	return fixPath, nil
}

// PagesFiles возвращает слайс со всеми page в виде абсолютных путей.
func PagesFiles(cfg config.StartupConfig) ([]string, error) {
	const operation = "helperFunc.PagesFiles"
	// Открываем директорию, чтобы прочитать все файлы
	d, err := os.Open(cfg.PagesPath)
	if err != nil {
		return nil, fmt.Errorf("%s: Не удалось открыть файловую директорию :%w", operation, err)

	}

	err = d.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: Не удалось закрыть файловую директорию :%w", operation, err)
	}

	// Получаем список файлов в директории
	files, err := d.Readdirnames(-1)
	if err != nil {
		return nil, fmt.Errorf("%s: Не удалось закрыть файловую директорию :%w", operation, err)
	}

	//Создаём абсолютные пути ко всем страницам
	for i, val := range files {
		files[i] = cfg.PagesPath + val
	}

	return files, nil
}
