package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Help     bool
	Parents  bool
	Verbose  bool
	Version  bool
	Dirs     []string
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

	if len(config.Dirs) == 0 {
		fmt.Fprintln(os.Stderr, "rmdir: пропущена директория")
		fmt.Fprintln(os.Stderr, "По команде «rmdir -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	executeRmdir(config)
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
		case "-p":
			config.Parents = true
			i++
			continue
		case "-v":
			config.Verbose = true
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
					case 'p':
						config.Parents = true
					case 'v':
						config.Verbose = true
					default:
						panic(fmt.Sprintf("rmdir: неверный ключ — '%s'", arg))
					}
				}
				i++
				continue
			}
			config.Dirs = append(config.Dirs, arg)
			i++
		}
	}

	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("rmdir - удаляет пустые директории")
	fmt.Println()
	fmt.Println("Использование: rmdir [ОПЦИЯ]... ДИРЕКТОРИИ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -p     удалять родительские директории, если они пустые")
	fmt.Println("  -v     подробный вывод")
	fmt.Println("  -h     показать эту справку")
	fmt.Println("  --version показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  rmdir dir1 dir2")
	fmt.Println("  rmdir -p -v path/to/dir")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("rmdir версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeRmdir удаляет директории согласно конфигурации
func executeRmdir(config *Config) {
	for _, dir := range config.Dirs {
		if err := removeDir(dir, config.Parents, config.Verbose); err != nil {
			fmt.Fprintf(os.Stderr, "rmdir: %v\n", err)
			os.Exit(1)
		}
	}
}

// removeDir удаляет директорию рекурсивно при необходимости
func removeDir(path string, parents bool, verbose bool) error {
	if parents {
		return removeParentDirs(path, verbose)
	}

	empty, err := isEmptyDir(path)
	if err != nil {
		return fmt.Errorf("ошибка доступа к '%s': %v", path, err)
	}
	if !empty {
		return fmt.Errorf("директория '%s' не пуста", path)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("не удалось удалить '%s': %v", path, err)
	}
	
	if verbose {
		fmt.Printf("Директория '%s' удалена.\n", path)
	}
	return nil
}

// isEmptyDir проверяет, пуста ли директория
func isEmptyDir(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	names, err := f.Readdirnames(0)
	if err != nil {
		return false, err
	}
	return len(names) == 0, nil
}

// removeParentDirs рекурсивно удаляет родительские пустые директории
func removeParentDirs(path string, verbose bool) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("не удалось удалить '%s': %v", path, err)
	}
	
	if verbose {
		fmt.Printf("Директория '%s' удалена.\n", path)
	}

	parent := filepath.Dir(path)
	if parent == "/" || parent == "." {
		return nil 
	}

	empty, err := isEmptyDir(parent)
	if err != nil {
		return nil
	}
	if empty {
		return removeParentDirs(parent, verbose) 
	}
	return nil
}

