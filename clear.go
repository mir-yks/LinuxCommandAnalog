package main

import (
	"flag"
	"fmt"
	"os"
)

type ClearOptions struct {
	Help      bool
	Version   bool
	CursesVer bool
	Force     bool
	NoScroll  bool
}

const (
	clearScreenCode = "\033[H\033[2J"
	clearLineCode   = "\033[2K"
	cursorHome      = "\033[H"
	clearScrollback = "\033[3J" 
)

func parseFlags() ClearOptions {
	var options ClearOptions
	
	flag.BoolVar(&options.Help, "h", false, "показать справку и выйти")
	flag.BoolVar(&options.Help, "help", false, "показать справку и выйти")
	flag.BoolVar(&options.Version, "v", false, "показать версию и выйти")
	flag.BoolVar(&options.Version, "version", false, "показать версию и выйти")
	flag.BoolVar(&options.CursesVer, "V", false, "показать версию curses и выйти")
	flag.BoolVar(&options.Force, "f", false, "принудительная очистка (игнорировать проверки)")
	flag.BoolVar(&options.Force, "force", false, "принудительная очистка (игнорировать проверки)")
	flag.BoolVar(&options.NoScroll, "x", false, "не очищать буфер прокрутки")
	
	flag.Parse()
	
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "clear: лишний операнд '%s'\n", flag.Arg(0))
		fmt.Fprintln(os.Stderr, "Используйте 'clear -h' для получения справки")
		os.Exit(1)
	}
	
	return options
}

func showHelp() {
	fmt.Println("Использование: clear [КЛЮЧ]...")
	fmt.Println()
	fmt.Println("Очищает экран терминала.")
	fmt.Println()
	fmt.Println("Ключи:")
	fmt.Println("  -h, --help     показать эту справку и выйти")
	fmt.Println("  -v, --version  показать информацию о версии и выйти")
	fmt.Println("  -V             показать версию curses и выйти")
	fmt.Println("  -f, --force    принудительная очистка (игнорировать проверки)")
	fmt.Println("  -x             не очищать буфер прокрутки")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  clear          # Очистить экран")
	fmt.Println("  clear -h       # Показать справку")
	fmt.Println("  clear -V       # Показать версию curses")
	fmt.Println("  clear -x       # Очистить только экран, не scrollback")
}

func showVersion() {
	fmt.Println("clear версия 1.0")
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

func showCursesVersion() {
	fmt.Println("ncurses 6.4.20230624") 
}

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

func clearScreen(options ClearOptions) error {
	if !options.Force && !isTerminalANSICapable() {
		term := os.Getenv("TERM")
		fmt.Fprintf(os.Stderr, "clear: терминал '%s' может не поддерживать очистку экрана\n", term)
		fmt.Fprintln(os.Stderr, "Используйте 'clear -f' для принудительной очистки")
		return fmt.Errorf("терминал не поддерживает ANSI escape коды")
	}
	
	if !options.NoScroll {
		fmt.Print(clearScrollback)
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
	
	if options.CursesVer {
		showCursesVersion()
		os.Exit(0)
	}
	
	if err := clearScreen(options); err != nil {
		fmt.Fprintf(os.Stderr, "clear: ошибка очистки экрана: %v\n", err)
		os.Exit(1)
	}
}

