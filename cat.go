package main

import (
	"bufio"
	"fmt"
	"os"
)

type Config struct {
	Help      bool
	Version   bool
	ShowAll   bool   
	NumNonEmpty bool 
	NumAll     bool    
	AddDollar  bool  
	Filenames  []string
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
		fmt.Fprintln(os.Stderr, "cat: пропущен операнд, задающий файл")
		fmt.Fprintln(os.Stderr, "По команде «cat -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	err := executeCat(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
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
		case "-A":
			config.ShowAll = true
		case "-b":
			config.NumNonEmpty = true
		case "-n":
			config.NumAll = true
		case "-e":
			config.AddDollar = true
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'v':
						config.Version = true
					case 'A':
						config.ShowAll = true
					case 'b':
						config.NumNonEmpty = true
					case 'n':
						config.NumAll = true
					case 'e':
						config.AddDollar = true
					default:
						panic(fmt.Sprintf("cat: неверный ключ — '%s'", arg))
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
	fmt.Println("cat - объединяет файлы и выводит их на стандартный вывод")
	fmt.Println()
	fmt.Println("Использование: cat [ОПЦИЯ]... ФАЙЛ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -A        эквивалентно -vET")
	fmt.Println("  -b        нумеровать непустые строки вывода")
	fmt.Println("  -e        заканчивать строки символом $")
	fmt.Println("  -n        нумеровать все строки вывода")
	fmt.Println("  -h        показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  cat file.txt             # Вывод содержимого файла")
	fmt.Println("  cat -n file.txt          # Нумеровать все строки")
	fmt.Println("  cat -b file.txt          # Нумеровать непустые строки")
	fmt.Println("  cat -e file.txt          # Показать символы конца строк")
	fmt.Println("  cat -A file.txt          # Показать все управляющие символы")
	fmt.Println("  cat file1.txt file2.txt  # Объединить несколько файлов")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("cat версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeCat выполняет основную логику команды cat
func executeCat(config *Config) error {
	for _, filename := range config.Filenames {
		err := readFile(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat: ошибка при обработке файла '%s': %v\n", filename, err)
		}
	}
	return nil
}

// readFile читает и обрабатывает файл согласно опциям
func readFile(fn string, config *Config) error {
	f, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		
		if config.NumAll || (config.NumNonEmpty && line != "") {
			fmt.Printf("%d\t", lineNum)
			lineNum++
		}

		if config.AddDollar {
			fmt.Println(line + "$")
		} else {
			fmt.Println(line)
		}

		if config.ShowAll {
			fmt.Printf("[showAll] %s\n", line)
		}
	}

	return scanner.Err()
}
