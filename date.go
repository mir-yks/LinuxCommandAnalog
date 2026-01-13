package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Help      bool
	Version   bool
	Date      string   
	File      string   
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

	err := executeDate(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
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
		case "-d":
			i++
			if i >= len(os.Args) {
				panic("touch: опция требует операнд -- 'd'")
			}
			config.Date = os.Args[i]
			i++
			continue
		case "-f", "-r":
			i++
			if i >= len(os.Args) {
				panic(fmt.Sprintf("date: опция '%s' требует операнд -- '%s'", arg, arg[1:]))
			}
			config.File = os.Args[i]
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'v':
						config.Version = true
					case 'd', 'f', 'r':
						panic(fmt.Sprintf("date: опция '%c' требует операнд", ch))
					default:
						panic(fmt.Sprintf("date: неверный ключ — '%s'", arg))
					}
				}
			} else {
				config.Filenames = append(config.Filenames, arg)
			}
			i++
		}
	}

	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("date - печатает или устанавливает системное время и дату")
	fmt.Println()
	fmt.Println("Использование: date [ОПЦИЯ]...") 
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -d <дата>    Печатает указанную дату (формат RFC3339)")
	fmt.Println("  -f <файл>    Читает даты из файла и печатает их")
	fmt.Println("  -r <файл>    Печатает время модификации файла")
	fmt.Println("  -h           Показать эту справку")
	fmt.Println("  -v, --version Показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  date                    # Текущее время")
	fmt.Println("  date -d '2026-01-13T12:00:00Z'  # Указанная дата")
	fmt.Println("  date -r file.txt         # Время файла")
	fmt.Println("  date -f dates.txt        # Даты из файла")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("date версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeDate выполняет основную логику команды date
func executeDate(config *Config) error {
	if config.Date == "" && config.File == "" && len(config.Filenames) == 0 {
		fmt.Println(time.Now().Format(time.RFC1123))
		return nil
	}

	// Обработка -d <дата>
	if config.Date != "" {
		if err := processDate(config.Date); err != nil {
			return fmt.Errorf("ошибка обработки даты: %v", err)
		}
	}

	// Обработка -f или -r <файл>
	if config.File != "" {
		if err := processFile(config.File, config); err != nil {
			return fmt.Errorf("ошибка обработки файла: %v", err)
		}
	}

	return nil
}

// processDate парсит и выводит дату
func processDate(dateInput string) error {
	if parsedDate, err := time.Parse(time.RFC3339, dateInput); err == nil {
		fmt.Println("Дата:", parsedDate.Format(time.RFC1123))
		return nil
	}
	return fmt.Errorf("неверный формат даты '%s'. Используйте RFC3339 (пример: 2026-01-13T12:00:00Z)", dateInput)
}

// processFile обрабатывает файл с датами или время файла
func processFile(filePath string, config *Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл '%s': %v", filePath, err)
	}
	defer file.Close()

	if config.File != "" && os.Args[1] == "-r" {
		info, err := file.Stat()
		if err != nil {
			return fmt.Errorf("не удалось получить информацию о файле: %v", err)
		}
		fmt.Println("Время файла:", info.ModTime().Format(time.RFC1123))
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if err := processDate(scanner.Text()); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка в строке: %v\n", err)
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("ошибка чтения файла: %v", err)
		}
	}

	return nil
}
