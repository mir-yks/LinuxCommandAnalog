package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Help      bool
	Version   bool
	Bytes     bool  
	Lines     bool  
	Words     bool  
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
		fmt.Fprintln(os.Stderr, "wc: пропущен операнд, задающий файл")
		fmt.Fprintln(os.Stderr, "По команде «wc -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	executeWc(config)
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
		case "-c":
			config.Bytes = true
		case "-l":
			config.Lines = true
		case "-w":
			config.Words = true
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'v':
						config.Version = true
					case 'c':
						config.Bytes = true
					case 'l':
						config.Lines = true
					case 'w':
						config.Words = true
					default:
						panic(fmt.Sprintf("wc: неверный ключ — '%s'", arg))
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
	fmt.Println("wc - подсчитывает количество строк, слов и байтов")
	fmt.Println()
	fmt.Println("Использование: wc [ОПЦИЯ]... ФАЙЛ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -c     количество байтов")
	fmt.Println("  -l     количество строк")
	fmt.Println("  -w     количество слов")
	fmt.Println("  -h     показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("По умолчанию выводятся все счетчики.")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  wc file.txt              # Все счетчики")
	fmt.Println("  wc -l file.txt           # Только строки")
	fmt.Println("  wc -w file.txt           # Только слова")
	fmt.Println("  wc -c file.txt           # Только байты")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("wc версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeWc выполняет основную логику команды wc
func executeWc(config *Config) {
	if !config.Bytes && !config.Lines && !config.Words {
		config.Bytes = true
		config.Lines = true
		config.Words = true
	}

	for _, filename := range config.Filenames {
		stats := countFile(filename)
		if stats != nil {
			printStats(filename, stats, config)
		}
	}
}

// countFile подсчитывает статистику файла
func countFile(filename string) *FileStats {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wc: %s: %v\n", filename, err)
		return nil
	}
	defer file.Close()

	var stats FileStats
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		stats.Bytes += len(line) + 1
		stats.Words += len(strings.Fields(line))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "wc: %s: %v\n", filename, err)
		return nil
	}

	return &stats
}

// printStats выводит статистику в формате wc
func printStats(filename string, stats *FileStats, config *Config) {
	if config.Bytes {
		fmt.Printf("%8d ", stats.Bytes)
	}
	if config.Lines {
		fmt.Printf("%8d ", stats.Lines)
	}
	if config.Words {
		fmt.Printf("%8d ", stats.Words)
	}
	fmt.Printf("%s\n", filename)
}

// FileStats структура для хранения статистики файла
type FileStats struct {
	Bytes  int
	Lines  int
	Words  int
}
