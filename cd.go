package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	Help     bool
	Version  bool
	Physical bool
	Logical  bool
	Verbose  bool
	Previous bool
	Path     string
}

const ver = "1.0.0"

func parseArgs() *Config {
	config := &Config{}
	
	// Создаем флаги для всех опций
	help := flag.Bool("h", false, "показать справку и выйти")
	helpLong := flag.Bool("help", false, "показать справку и выйти")
	version := flag.Bool("version", false, "показать информацию о версии и выйти")
	physical := flag.Bool("P", false, "физическая смена директории")
	logical := flag.Bool("L", false, "логическая смена директории")
	verbose := flag.Bool("v", false, "выводить подробную информацию")
	
	flag.Parse()
	
	config.Help = *help || *helpLong
	config.Version = *version
	config.Physical = *physical
	config.Logical = *logical
	config.Verbose = *verbose
	
	// Проверяем позиционные аргументы
	switch flag.NArg() {
	case 0:
		// Путь не указан
	case 1:
		if flag.Arg(0) == "-" {
			config.Previous = true
		} else {
			config.Path = flag.Arg(0)
		}
	default:
		fmt.Fprintf(os.Stderr, "cd: лишний операнд '%s'\n", flag.Arg(1))
		fmt.Fprintln(os.Stderr, "Используйте 'cd -h' для получения справки")
		os.Exit(1)
	}
	
	return config
}

func resolvePath(path string) string {
	absPath, _ := filepath.Abs(path)
	return absPath
}

func executeCd(config *Config) {
	var finalPath string

	if config.Previous {
		current, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: ошибка получения текущей директории: %v\n", err)
			os.Exit(1)
		}
		finalPath = filepath.Dir(current)
	} else {
		if config.Path == "" {
			fmt.Fprintln(os.Stderr, "cd: не указан путь")
			fmt.Fprintln(os.Stderr, "Используйте 'cd -h' для получения справки")
			os.Exit(1)
		}
		finalPath = resolvePath(config.Path)
	}

	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "cd: '%s': нет такого файла или директории\n", finalPath)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Printf("Переход в: %s\n", finalPath)
		fmt.Printf("Абсолютный путь: %s\n", finalPath)
		current, _ := os.Getwd()
		fmt.Printf("Было: %s\n", current)
		if config.Physical {
			fmt.Println("Режим: физический")
		} else if config.Logical {
			fmt.Println("Режим: логический")
		}
	}

	cmd := exec.Command("bash")
	cmd.Dir = finalPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cd: ошибка запуска shell: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Использование: cd [КЛЮЧ]... [ПУТЬ|-]")
	fmt.Println()
	fmt.Println("Переходит в указанную директорию и запускает bash shell.")
	fmt.Println()
	fmt.Println("Ключи:")
	fmt.Println("  -h, --help        показать эту справку и выйти")
	fmt.Println("  --version         показать информацию о версии и выйти")
	fmt.Println("  -P                физическая смена директории")
	fmt.Println("  -L                логическая смена директории")
	fmt.Println("  -v                выводить подробную информацию")
	fmt.Println("  -                 перейти в предыдущую директорию")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  cd -                  # Перейти в предыдущую директорию")
	fmt.Println("  cd -v tmp/            # Подробная информация")
	fmt.Println("  cd -P tmp/            # Физическая смена директории")
	fmt.Println("  cd -L tmp/            # Логическая смена директории")
	fmt.Println("  cd -h                 # Показать справку")
}

func printVersion() {
	fmt.Println("cd версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

func main() {
	config := parseArgs()

	if config.Help {
		printHelp()
		os.Exit(0)
	}

	if config.Version {
		printVersion()
		os.Exit(0)
	}

	executeCd(config)
}

