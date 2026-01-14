package main

import (
	"fmt"
	"os"
)

type Config struct {
	Help      bool   
	Version   bool   
	Force     bool   
	Recursive bool   
	Verbose   bool   
	Paths     []string 
}

const ver = "1.0.0"

// printHelp выводит справку по использованию утилиты rm
func printHelp() {
	fmt.Println("rm - удаление файлов и директорий")
	fmt.Println()
	fmt.Println("Использование: rm [ОПЦИЯ]... ПУТЬ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -h, --help     Показать эту справку")
	fmt.Println("  -v             Подробный вывод")
	fmt.Println("  --version      Показать информацию о версии")
	fmt.Println("  -f             Принудительное удаление (игнорировать ошибки)")
	fmt.Println("  -R, -r         Рекурсивное удаление директорий")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  rm file.txt                 # Удалить файл")
	fmt.Println("  rm -v file.txt              # С выводом")
	fmt.Println("  rm -f file.txt              # Принудительно")
	fmt.Println("  rm -R ./dir                  # Рекурсивно")
}

// printVersion выводит информацию о версии программы
func printVersion() {
	fmt.Println("rm версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// parseArgs разбирает аргументы командной строки в стиле ls
func parseArgs() *Config {
	config := &Config{}
	i := 1

	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-h", "--help":
			config.Help = true
			i++
			continue
		case "-v", "--version":
			if arg == "-v" {
				config.Verbose = true
			} else {
				config.Version = true
			}
			i++
			continue
		case "-f":
			config.Force = true
			i++
			continue
		case "-R", "-r":
			config.Recursive = true
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				panic(fmt.Sprintf("rm: неверный ключ '%s'. Используйте --help для справки", arg))
			}
			config.Paths = append(config.Paths, arg)
			i++
		}
	}

	return config
}

// removeFile удаляет файл или директорию
func removeFile(path string, force, recursive, verbose bool) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: нет такого файла или директории", path)
	}

	if info.IsDir() && !recursive {
		return fmt.Errorf("%s: это директория. Используйте -R или -r", path)
	}

	if recursive {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}

	if verbose && err == nil {
		fmt.Printf("удален: %s\n", path)
	}

	return err
}

// executeRm выполняет удаление файлов/директорий
func executeRm(config *Config) {
	if len(config.Paths) == 0 {
		panic("rm: пропущен операнд")
	}

	for _, path := range config.Paths {
		err := removeFile(path, config.Force, config.Recursive, config.Verbose)
		
		if err != nil {
			if config.Force {
				fmt.Fprintf(os.Stderr, "rm: %v\n", err)
			} else {
				panic(err)
			}
		}
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "rm: %v\n", r)
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

	executeRm(config)
}

