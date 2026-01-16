package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	ClearHist  bool
	DeleteLine int
	NumLines   int
	Help       bool
}

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

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		fmt.Fprintln(os.Stderr, "history: не задана переменная HOME")
		os.Exit(1)
	}
	histFile := filepath.Join(homeDir, ".bash_history")

	if config.ClearHist {
		err := os.Truncate(histFile, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "history: не удалось очистить %s: %v\n", histFile, err)
			os.Exit(1)
		}
		fmt.Println("История очищена")
		return
	}

	if config.DeleteLine > 0 {
		deleteLineFromHistory(histFile, config.DeleteLine)
		return
	}

	showHistory(histFile, config.NumLines)
}

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
		case "-c":
			config.ClearHist = true
			i++
			continue
		case "-d":
			i++
			if i >= len(os.Args) {
				panic("history: ожидается номер строки после -d")
			}
			num, err := strconv.Atoi(os.Args[i])
			if err != nil || num < 1 {
				panic("history: номер строки после -d должен быть положительным числом")
			}
			config.DeleteLine = num
			i++
			continue
		case "-n":
			i++
			if i >= len(os.Args) {
				panic("history: ожидается число после -n")
			}
			num, err := strconv.Atoi(os.Args[i])
			if err != nil || num < 1 {
				panic("history: число после -n должно быть положительным")
			}
			config.NumLines = num
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'c':
						config.ClearHist = true
					case 'd':
						panic("history: неверный формат -d (ожидается -d N)")
					case 'n':
						panic("history: неверный формат -n (ожидается -n N)")
					default:
						panic(fmt.Sprintf("history: неверный ключ — '%s'", arg))
					}
				}
				i++
				continue
			}
			panic(fmt.Sprintf("history: неизвестный аргумент '%s'", arg))
		}
	}

	return config
}

func printHelp() {
	fmt.Println("history - просмотр и управление историей команд bash")
	fmt.Println()
	fmt.Println("Использование: history [ОПЦИЯ]...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -c        очистить файл истории")
	fmt.Println("  -d N      удалить строку N из истории")
	fmt.Println("  -n N      показать последние N строк")
	fmt.Println("  -h        показать эту справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  history                           # Показать всю историю")
	fmt.Println("  history -n 10                     # Последние 10 команд")
	fmt.Println("  history -d 5                      # Удалить 5-ю строку")
	fmt.Println("  history -c                        # Очистить историю")
}

func showHistory(histFile string, numLines int) {
	file, err := os.Open(histFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Fprintf(os.Stderr, "history: не удалось открыть %s: %v\n", histFile, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "history: ошибка чтения %s: %v\n", histFile, err)
		os.Exit(1)
	}

	if numLines > 0 && len(lines) > numLines {
		lines = lines[len(lines)-numLines:]
	}

	for i, line := range lines {
		fmt.Printf("%5d  %s\n", i+1, line)
	}
}

func deleteLineFromHistory(histFile string, lineNum int) {
	lines := readLines(histFile)
	if lineNum < 1 || lineNum > len(lines) {
		fmt.Fprintf(os.Stderr, "history: строка %d вне диапазона (1–%d)\n", lineNum, len(lines))
		os.Exit(1)
	}

	lines = append(lines[:lineNum-1], lines[lineNum:]...)

	err := writeLines(histFile, lines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: не удалось записать %s: %v\n", histFile, err)
		os.Exit(1)
	}
	fmt.Printf("Удалена строка %d\n", lineNum)
}

func readLines(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func writeLines(file string, lines []string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		_, err := fmt.Fprintln(f, line)
		if err != nil {
			return err
		}
	}
	return f.Sync()
}
