package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Help     bool  
	Version  bool   
	Summary  bool   
	All      bool   
	Path     string 
}

const ver = "1.0.0"

// printHelp выводит справку по использованию утилиты du
func printHelp() {
	fmt.Println("du - оценка использования дискового пространства")
	fmt.Println()
	fmt.Println("Использование: du [ОПЦИЯ]... [ПУТЬ]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -h, --help     Показать эту справку")
	fmt.Println("  -v, --version  Показать информацию о версии")
	fmt.Println("  -s             Вывести только общий размер")
	fmt.Println("  -a             Показать размер каждого файла")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  du                     # Размер текущей директории")
	fmt.Println("  du -s /tmp            # Общий размер /tmp")
	fmt.Println("  du -a dir/            # Каждый файл отдельно")
}

// printVersion выводит информацию о версии программы
func printVersion() {
	fmt.Println("du версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// parseArgs разбирает аргументы командной строки в стиле ls/rm
func parseArgs() *Config {
	config := &Config{Path: "."}
	i := 1

	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-h", "--help":
			config.Help = true
			i++
			continue
		case "-v", "--version":
			config.Version = true
			i++
			continue
		case "-s":
			config.Summary = true
			i++
			continue
		case "-a":
			config.All = true
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				panic(fmt.Sprintf("du: неверный ключ '%s'. Используйте --help для справки", arg))
			}
			config.Path = arg
			i++
		}
	}

	return config
}

// getSizes вычисляет размер каждого файла и возвращает map
func getSizes(path string, all bool) (map[string]uint64, uint64, error) {
	sizes := make(map[string]uint64)
	var totalSize uint64

	err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		fileSize := uint64(info.Size())
		
		totalSize += fileSize
		
		if all && !info.IsDir() {
			relPath, _ := filepath.Rel(path, p)
			sizes[relPath] = fileSize
		}

		return nil
	})

	return sizes, totalSize, err
}

// executeDu выполняет основную логику утилиты du
func executeDu(config *Config) {
	sizes, totalSize, err := getSizes(config.Path, config.All)
	if err != nil {
		panic(fmt.Sprintf("du: не удается прочитать '%s': %v", config.Path, err))
	}

	if config.All {
		for relPath, size := range sizes {
			fmt.Printf("%d\t%s\n", size, relPath)
		}
		fmt.Printf("%d\t%s\n", totalSize, config.Path)
	} else if config.Summary {
		fmt.Printf("%d\t%s\n", totalSize, config.Path)
	} else {
		fmt.Printf("Размер '%s': %d байт\n", config.Path, totalSize)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "du: %v\n", r)
			os.Exit(1)
		}
	}()

	config := parseArgs()

	if config.Help {
		printHelp()
		return
	}

	if config.Version {
		printVersion()
		return
	}

	executeDu(config)
}

