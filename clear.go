package main

import (
	"flag"
	"fmt"
	"os"
)

// Структура для хранения параметров команды
type ClearOptions struct {
	Help    bool   
	Version bool   
	Force   bool   
}

// ANSI escape коды для управления терминалом
const (
	clearScreenCode = "\033[H\033[2J"  
	clearLineCode   = "\033[2K"        
	cursorHome      = "\033[H"         
)

func parseFlags() ClearOptions {
	var options ClearOptions
	
	flag.BoolVar(&options.Help, "h", false, "показать справку и выйти")
	flag.BoolVar(&options.Help, "help", false, "показать справку и выйти")
	flag.BoolVar(&options.Version, "v", false, "показать версию и выйти")
	flag.BoolVar(&options.Version, "version", false, "показать версию и выйти")
	flag.BoolVar(&options.Force, "f", false, "принудительная очистка (игнорировать проверки)")
	flag.BoolVar(&options.Force, "force", false, "принудительная очистка (игнорировать проверки)")
	
	flag.Parse()
	
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "clear: лишний операнд '%s'\n", flag.Arg(0))
		fmt.Fprintln(os.Stderr, "Используйте 'clear -h' для получения справки")
		os.Exit(1)
	}
	
	return options
}

// Функция для отображения справки
func showHelp() {
	fmt.Println("Использование: clear [КЛЮЧ]...")
	fmt.Println()
	fmt.Println("Очищает экран терминала.")
	fmt.Println()
	fmt.Println("Ключи:")
	fmt.Println("  -h, --help     показать эту справку и выйти")
	fmt.Println("  -v, --version  показать информацию о версии и выйти")
	fmt.Println("  -f, --force    принудительная очистка (игнорировать проверки терминала)")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  clear          # Очистить экран")
	fmt.Println("  clear -h       # Показать справку")
	fmt.Println("  clear -v       # Показать версию")
	fmt.Println("  clear -f       # Принудительная очистка")
}

// Функция для отображения версии
func showVersion() {
	fmt.Println("clear версия 1.0")
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// Функция для проверки поддержки ANSI escape кодов
func isTerminalANSICapable() bool {
	term := os.Getenv("TERM")
	
	ansiTerms := []string{
		"xterm", "xterm-256color", "xterm-color",
		"screen", "screen-256color",
		"linux", "vt100", "vt220",
		"rxvt", "rxvt-unicode",
		"eterm", "ansi",
	}
	
	for _, ansiTerm := range ansiTerms {
		if term == ansiTerm {
			return true
		}
	}
	
	return false
}

// Функция очистки экрана
func clearScreen(options ClearOptions) error {
	if !options.Force && !isTerminalANSICapable() {
		fmt.Fprintln(os.Stderr, "clear: терминал может не поддерживать очистку экрана")
		fmt.Fprintln(os.Stderr, "Используйте 'clear -f' для принудительной очистки")
		return fmt.Errorf("терминал не поддерживает ANSI escape коды")
	}
	
	fmt.Print(clearScreenCode)
	
	fmt.Print(cursorHome)
	
	os.Stdout.Sync()
	
	return nil
}

func main() {
	options := parseFlags()
	
	if options.Help {
		showHelp()
		os.Exit(0)
	}
	
	if options.Version {
		showVersion()
		os.Exit(0)
	}
	
	if err := clearScreen(options); err != nil {
		fmt.Fprintf(os.Stderr, "clear: ошибка очистки экрана: %v\n", err)
		os.Exit(1)
	}
}
