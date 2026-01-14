package main

import (
	"fmt"
	"os"
)

type Config struct {
	Help        bool
	Parents     bool
	Verbose     bool
	Version     bool
	Directories []string
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

	if len(config.Directories) == 0 {
		fmt.Fprintln(os.Stderr, "mkdir: пропущен путь")
		fmt.Fprintln(os.Stderr, "По команде «mkdir -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	executeMkdir(config)
}

// parseArgs разбирает аргументы командной строки вручную
func parseArgs() *Config {
	config := &Config{}
	i := 1

	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-h":
			config.Help = true
			i++
			continue
		case "-v":
			config.Verbose = true
			i++
			continue
		case "-p":
			config.Parents = true
			i++
			continue
		case "--version":
			config.Version = true
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'v':
						config.Verbose = true
					case 'p':
						config.Parents = true
					default:
						panic(fmt.Sprintf("mkdir: неверный ключ — '%s'", arg))
					}
				}
				i++
				continue
			}
			config.Directories = append(config.Directories, arg)
			i++
		}
	}

	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("mkdir - создает директории")
	fmt.Println()
	fmt.Println("Использование: mkdir [ОПЦИЯ]... ДИРЕКТОРИИ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -p     создать родительские директории по необходимости")
	fmt.Println("  -v     подробный вывод")
	fmt.Println("  -h     показать эту справку")
	fmt.Println("  --version показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  mkdir dir1 dir2")
	fmt.Println("  mkdir -p ./path/to/new/dir")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("mkdir версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeMkdir создает директории согласно конфигурации
func executeMkdir(config *Config) {
	for _, dir := range config.Directories {
		if err := createDir(dir, config.Parents, config.Verbose); err != nil {
			fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
			os.Exit(1)
		}
	}
}

// createDir создает директорию с заданными параметрами
func createDir(path string, parents bool, verbose bool) error {
	mode := os.FileMode(0755)
	if parents {
		err := os.MkdirAll(path, mode)
		if err != nil {
			return fmt.Errorf("не удалось создать '%s': %v", path, err)
		}
		if verbose {
			fmt.Printf("Директория '%s' создана.\n", path)
		}
	} else {
		err := os.Mkdir(path, mode)
		if err != nil {
			return fmt.Errorf("не удалось создать '%s': %v", path, err)
		}
		if verbose {
			fmt.Printf("Директория '%s' создана.\n", path)
		}
	}
	return nil
}

