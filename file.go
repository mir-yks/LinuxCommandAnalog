package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Help     bool
	Version  bool
	Filenames []string
}

const ver = "1.0.0"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", r)
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

	if len(config.Filenames) == 0 {
		fmt.Fprintln(os.Stderr, "file: пропущен операнд, задающий файл")
		fmt.Fprintln(os.Stderr, "По команде «file -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	executeFile(config)
}

// parseArgs разбирает аргументы командной строки вручную
func parseArgs() *Config {
	config := &Config{}
	filenames := []string{}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-h":
			config.Help = true
		case "-v", "--version":
			config.Version = true
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'v':
						config.Version = true
					default:
						panic(fmt.Sprintf("file: неверный ключ — '%s'", arg))
					}
				}
			} else {
				filenames = append(filenames, arg)
			}
		}
	}

	config.Filenames = filenames
	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("file - определяет тип файла")
	fmt.Println()
	fmt.Println("Использование: file [ОПЦИЯ]... ФАЙЛ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -h       показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("Определяет тип файла по расширению.")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  file document.txt      # Текстовый файл")
	fmt.Println("  file image.jpg         # Изображение JPEG")
	fmt.Println("  file program.go        # Файл Go")
	fmt.Println("  file archive.zip       # ZIP архив")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("file версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeFile выполняет основную логику команды file
func executeFile(config *Config) {
	for _, filename := range config.Filenames {
		if info, err := os.Stat(filename); err != nil || info.IsDir() {
			fmt.Fprintf(os.Stderr, "file: %s: не файл\n", filename)
			continue
		}
		filetype := getFileType(filename)
		fmt.Printf("%s: %s\n", filename, filetype)
	}
}

// getFileType определяет тип файла по расширению
func getFileType(name string) string {
	ext := filepath.Ext(name)

	switch ext {
	case ".txt":
		return "текстовый файл"
	case ".jpg", ".jpeg":
		return "изображение JPEG"
	case ".png":
		return "изображение PNG"
	case ".gif":
		return "изображение GIF"
	case ".go":
		return "файл Go"
	case ".c":
		return "файл C"
	case ".py":
		return "файл Python"
	case ".pdf":
		return "PDF документ"
	case ".zip":
		return "ZIP архив"
	case ".tar", ".gz":
		return "архив"
	case ".sh":
		return "shell-скрипт"
	case ".html", ".htm":
		return "HTML документ"
	case ".css":
		return "CSS файл"
	case ".json":
		return "JSON файл"
	default:
		return "неизвестный тип файла"
	}
}
