package main

import (
	"fmt"
	"os"
	"runtime"
)

type Config struct {
	Help     bool
	Version  bool
	All      bool
	Kernel   bool
	Host     bool
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

	if !config.All && !config.Kernel && !config.Host {
		config.Kernel = true 
	}

	executeUname(config)
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
		case "-v", "--version":
			config.Version = true
			i++
			continue
		case "-a":
			config.All = true
			i++
			continue
		case "-s":
			config.Kernel = true
			i++
			continue
		case "-n":
			config.Host = true
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'a':
						config.All = true
					case 's':
						config.Kernel = true
					case 'n':
						config.Host = true
					case 'v':
						config.Version = true
					case 'h':
						config.Help = true
					default:
						panic(fmt.Sprintf("uname: неверный ключ — '%s'", arg))
					}
				}
				i++
				continue
			}
			panic(fmt.Sprintf("неожиданный аргумент: %s", arg))
		}
	}

	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("uname - выводит информацию о системе")
	fmt.Println()
	fmt.Println("Использование: uname [ОПЦИЯ]...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -a     показать всю информацию о системе")
	fmt.Println("  -s     показать имя ядра (по умолчанию)")
	fmt.Println("  -n     показать имя узла")
	fmt.Println("  -h     показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  uname -a                    # Вся информация")
	fmt.Println("  uname                      # Имя ядра (по умолчанию)")
	fmt.Println("  uname -n                   # Имя узла")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("uname версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeUname выполняет вывод информации о системе
func executeUname(config *Config) {
	sysInfo := getSystemInfo()

	if config.All {
		sysInfo.Show()
		return
	}

	if config.Kernel {
		fmt.Println(sysInfo.Kernel)
	}
	if config.Host {
		fmt.Println(sysInfo.Host)
	}
}

// getSystemInfo собирает информацию о системе
type SystemInfo struct {
	Kernel string
	Host   string
}

func getSystemInfo() *SystemInfo {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	return &SystemInfo{
		Kernel: runtime.GOOS,
		Host:   hostname,
	}
}

func (si *SystemInfo) Show() {
	fmt.Printf("Имя ядра: %s\n", si.Kernel)
	fmt.Printf("Имя узла: %s\n", si.Host)
}

