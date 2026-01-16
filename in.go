package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Help    bool
	Version bool
	Number  int
}

const ver = "1.0.0"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "!n: %v\n", r)
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

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		fmt.Fprintln(os.Stderr, "!n: не задана переменная HOME")
		os.Exit(1)
	}
	histFile := filepath.Join(homeDir, ".bash_history")

	history := readHistoryFromFile(histFile)
	if len(history) == 0 {
		fmt.Fprintln(os.Stderr, "!n: история пуста")
		os.Exit(1)
	}

	if config.Number == 0 {
		fmt.Print("!n: введите номер команды: ")
		var input string
		fmt.Scanln(&input)
		num, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || num < 1 {
			fmt.Fprintln(os.Stderr, "!n: номер команды должен быть положительным числом")
			os.Exit(1)
		}
		config.Number = num
	}

	if config.Number > len(history) {
		fmt.Fprintf(os.Stderr, "!n: команда #%d не существует (всего %d)\n", config.Number, len(history))
		os.Exit(1)
	}

	cmdLine := history[config.Number-1]
	parts := strings.Fields(cmdLine)
	if len(parts) == 0 {
		fmt.Fprintln(os.Stderr, "!n: пустая команда в истории")
		os.Exit(1)
	}

	finalArgs := parts
	userArgs := getUserArgs(config)
	finalArgs = append(finalArgs, userArgs...)

	fmt.Printf("%s\n", strings.Join(finalArgs, " "))

	c := exec.Command(finalArgs[0], finalArgs[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "!n: выполнение '%s' завершилось с ошибкой: %v\n", finalArgs[0], err)
		os.Exit(1)
	}
}

func getUserArgs(config *Config) []string {
	var args []string
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-h" || arg == "-v" || arg == "--version" {
			continue
		}
		if config.Number > 0 && isNumberArg(arg) {
			continue
		}
		args = append(args, arg)
	}
	return args
}

func isNumberArg(arg string) bool {
	num, err := strconv.Atoi(arg)
	return err == nil && num > 0
}

func parseArgs() *Config {
	config := &Config{}
	
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-h":
			config.Help = true
			return config
		case "-v", "--version":
			config.Version = true
			return config
		default:
			num, err := strconv.Atoi(arg)
			if err == nil && num > 0 {
				config.Number = num
				return config
			}
		}
	}
	return config
}

func printHelp() {
	fmt.Println("!n - выполнение команд из истории по номеру")
	fmt.Println()
	fmt.Println("Использование: !n [ОПЦИЯ]... [N] [аргументы команды]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -h      показать эту справку")
	fmt.Println("  -v, --version  показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  !n 92                       # Выполнить 92-ю команду (ls)")
	fmt.Println("  !n 92 -la                   # ls -la")
	fmt.Println("  !n                          # Запрос номера команды")
}

func printVersion() {
	fmt.Println("!n версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

func readHistoryFromFile(histFile string) []string {
	file, err := os.Open(histFile)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("ошибка чтения истории: %v", err))
	}
	return lines
}

