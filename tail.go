package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Config struct {
	Help     bool
	Version  bool
	Bytes    int      
	Lines    int      
	Filename string
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

	if config.Filename == "" {
		fmt.Fprintln(os.Stderr, "tail: пропущен операнд, задающий файл")
		fmt.Fprintln(os.Stderr, "По команде «tail -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	if config.Bytes > 0 {
		if err := readLastBytes(config.Filename, config.Bytes); err != nil {
			fmt.Fprintf(os.Stderr, "tail: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := readLastLines(config.Filename, config.Lines); err != nil {
		fmt.Fprintf(os.Stderr, "tail: %v\n", err)
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
		case "-c", "-n":
			i++
			if i >= len(os.Args) {
				panic(fmt.Sprintf("tail: опция '%s' требует операнд", arg))
			}
			count, err := strconv.Atoi(os.Args[i])
			if err != nil || count < 0 {
				panic(fmt.Sprintf("tail: неверное значение '%s' для опции %s", os.Args[i], arg))
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
					default:
						panic(fmt.Sprintf("tail: неверный ключ — '%s'", arg))
					}
				}
			} else {
				config.Filename = arg
			}
			i++
		}
	}

	if config.Lines == 0 && config.Bytes == 0 {
		config.Lines = 10
	}

	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("tail - выводит конец каждого заданного файла")
	fmt.Println()
	fmt.Println("Использование: tail [ОПЦИЯ]... ФАЙЛ")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -c N     выводить последние N байт")
	fmt.Println("  -n N     выводить последние N строк")
	fmt.Println("  -h       показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("По умолчанию N=10 строк.")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  tail file.txt              # Последние 10 строк")
	fmt.Println("  tail -n 5 file.txt         # Последние 5 строк")
	fmt.Println("  tail -c 100 file.txt       # Последние 100 байт")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("tail версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// readLastLines читает последние N строк из файла
func readLastLines(filePath string, lineCount int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть '%s': %v", filePath, err)
	}
	defer file.Close()

	// Проверяем размер файла
	stat, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о файле: %v", err)
	}
	
	// Если файл пустой
	if stat == 0 {
		return nil
	}
	
	// Возвращаемся в начало
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("ошибка перемещения: %v", err)
	}

	var allLines []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения: %v", err)
	}

	startIndex := 0
	if len(allLines) > lineCount {
		startIndex = len(allLines) - lineCount
	}
	
	for i := startIndex; i < len(allLines); i++ {
		fmt.Println(allLines[i])
	}
	
	return nil
}

// readLastBytes читает последние N байт из файла
func readLastBytes(filePath string, byteCount int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть '%s': %v", filePath, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о файле: %v", err)
	}

	start := int64(0)
	if stat.Size() > int64(byteCount) {
		start = stat.Size() - int64(byteCount)
	}

	if _, err := file.Seek(start, io.SeekStart); err != nil {
		return fmt.Errorf("ошибка перемещения: %v", err)
	}

	buffer := make([]byte, byteCount)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("ошибка чтения: %v", err)
	}

	os.Stdout.Write(buffer[:n])
	fmt.Println() 
	return nil
}

