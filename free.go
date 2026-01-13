package main

import (
	"fmt"
	"os"
	"syscall"
)

type Config struct {
	Help   bool
	Version bool
	Bytes  bool 
	Mega   bool   
	Giga   bool 
}

const ver = "1.0.0"

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

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

	executeFree(config)
}

// parseArgs разбирает аргументы командной строки вручную
func parseArgs() *Config {
	config := &Config{}
	
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-h":
			config.Help = true
		case "-v", "--version":
			config.Version = true
		case "-b":
			config.Bytes = true
		case "-m":
			config.Mega = true
		case "-g":
			config.Giga = true
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'v':
						config.Version = true
					case 'b':
						config.Bytes = true
					case 'm':
						config.Mega = true
					case 'g':
						config.Giga = true
					default:
						panic(fmt.Sprintf("free: неверный ключ — '%s'", arg))
					}
				}
			} else {
				panic(fmt.Sprintf("free: неизвестный аргумент '%s'", arg))
			}
		}
	}

	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("free - отображает информацию о памяти")
	fmt.Println()
	fmt.Println("Использование: free [ОПЦИЯ]...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -b     байты")
	fmt.Println("  -m     мегабайты")
	fmt.Println("  -g     гигабайты")
	fmt.Println("  -h     показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  free                    # По умолчанию в MB")
	fmt.Println("  free -m                 # В мегабайтах")
	fmt.Println("  free -g                 # В гигабайтах")
	fmt.Println("  free -b                 # В байтах")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("free версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeFree выполняет основную логику команды free
func executeFree(config *Config) {
	total, used, free := getMemoryStats()
	
	var unit uint64 = MB 
	if config.Bytes {
		unit = 1
	} else if config.Giga {
		unit = GB
	} else if config.Mega {
		unit = MB
	}

	printMemoryStats(total, used, free, unit)
}

// getMemoryStats получает статистику памяти
func getMemoryStats() (total, used, free uint64) {
	var sysInfo syscall.Sysinfo_t
	syscall.Sysinfo(&sysInfo)

	total = sysInfo.Totalram * uint64(sysInfo.Unit)
	free = sysInfo.Freeram * uint64(sysInfo.Unit)
	used = total - free
	return
}

// printMemoryStats выводит статистику памяти
func printMemoryStats(total, used, free uint64, unit uint64) {
	unitName := getUnitName(unit)
	
	fmt.Printf("Mem: %s total, %s used, %s free\n",
		formatBytes(total, unit), 
		formatBytes(used, unit), 
		formatBytes(free, unit))
	_ = unitName
}

// formatBytes форматирует количество байт
func formatBytes(bytes, unit uint64) string {
	value := bytes / unit
	return fmt.Sprintf("%d %s", value, getUnitName(unit))
}

// getUnitName возвращает название единицы измерения
func getUnitName(unit uint64) string {
	switch unit {
	case 1:
		return "B"
	case KB:
		return "KB"
	case MB:
		return "MB"
	case GB:
		return "GB"
	default:
		return "B"
	}
}

