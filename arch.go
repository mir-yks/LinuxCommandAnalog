package main

import (
	"fmt"
	"os"
	"runtime"
)

type Config struct {
	Help      bool 
	Version   bool 
	Verbose   bool 
}

const ver = "1.0.0"

// printHelp выводит справку по использованию утилиты arch
func printHelp() {
	fmt.Println("arch - вывод информации об архитектуре системы")
	fmt.Println()
	fmt.Println("Использование: arch [ОПЦИЯ]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -h, --help     Показать эту справку")
	fmt.Println("  -v             Подробный вывод")
	fmt.Println("  --version      Показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  arch                 # Архитектура")
	fmt.Println("  arch -v              # Подробная информация") 
	fmt.Println("  arch --version       # Версия программы")
}

// printVersion выводит информацию о версии программы
func printVersion() {
	fmt.Println("arch версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// parseArgs разбирает аргументы командной строки в стиле ls/zip
func parseArgs() *Config {
	config := &Config{}
	
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-h", "--help":
			config.Help = true
			continue
		case "-v":
			config.Verbose = true
			continue
		case "--version":
			config.Version = true
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				panic(fmt.Sprintf("arch: неизвестная опция '%s'. Используйте -h для справки", arg))
			}
		}
	}

	return config
}

// printArchitecture выводит информацию об архитектуре системы
func printArchitecture(verbose bool) {
	arch := runtime.GOARCH
	osName := runtime.GOOS
	
	fmt.Printf("Архитектура: %s\n", arch)
	
	if verbose {
		fmt.Printf("ОС: %s\n", osName)
		fmt.Printf("GOOS: %s, GOARCH: %s\n", osName, arch)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "arch: %v\n", r)
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

	printArchitecture(config.Verbose)
}

