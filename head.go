package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Help      bool
	Version   bool
	Bytes     int      
	Lines     int      
	Quiet     bool     
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
		fmt.Fprintln(os.Stderr, "head: пропущен операнд, задающий файл")
		fmt.Fprintln(os.Stderr, "По команде «head -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	err := executeHead(config)
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
		case "-q":
			config.Quiet = true
			i++
			continue
		case "-c", "-n":
			i++
			if i >= len(os.Args) {
				panic(fmt.Sprintf("head: опция '%s' требует операнд", arg))
			}
			count, err := strconv.Atoi(os.Args[i])
			if err != nil || count < 0 {
				panic(fmt.Sprintf("head: неверное значение '%s' для опции %s", os.Args[i], arg))
			}
			if arg == "-c" {
				config.Bytes = count
			} else {
				config.Lines = count
			}
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
					case 'q':
						config.Quiet = true
					default:
						panic(fmt.Sprintf("head: неверный ключ — '%s'", arg))
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
	fmt.Println("head - выводит начало каждого заданного файла")
	fmt.Println()
	fmt.Println("Использование: head [ОПЦИЯ]... ФАЙЛ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -c N     выводить первые N байт")
	fmt.Println("  -n N     выводить первые N строк")
	fmt.Println("  -q       никогда не выводить заголовки файлов")
	fmt.Println("  -h       показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("По умолчанию N=10 строк. Если ни -c, ни -n не указаны,")
	fmt.Println("выводятся первые 10 строк.")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  head file.txt              # Первые 10 строк")
	fmt.Println("  head -n 5 file.txt         # Первые 5 строк")
	fmt.Println("  head -c 100 file.txt       # Первые 100 байт")
	fmt.Println("  head -q file1.txt file2.txt # Без заголовков")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("head версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeHead выполняет основную логику команды head
func executeHead(config *Config) error {
	if config.Lines == 0 && config.Bytes == 0 {
		config.Lines = 10
	}

	for _, filename := range config.Filenames {
		err := readFile(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "head: ошибка при обработке файла '%s': %v\n", filename, err)
		}
	}
	return nil
}

// readFile читает и выводит начало файла
func readFile(path string, config *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("не удалось открыть '%s': %v", path, err)
	}
	defer file.Close()

	if config.Bytes > 0 {
		if !config.Quiet {
			fmt.Printf("%s\n", path)
		}
		buffer := make([]byte, config.Bytes)
		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			return fmt.Errorf("ошибка чтения: %v", err)
		}
		os.Stdout.Write(buffer[:n])
		if !config.Quiet {
			fmt.Println()
		}
	} else {
		if !config.Quiet {
			fmt.Printf("%s\n", path)
		}
		scanner := bufio.NewScanner(file)
		count := 0
		for count < config.Lines && scanner.Scan() {
			fmt.Println(scanner.Text())
			count++
		}
	}

	return nil
}

