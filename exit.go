package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

func main() {
	// Определяем флаги
	help := flag.Bool("h", false, "Показать помощь и выйти")
	helpLong := flag.Bool("help", false, "Показать помощь и выйти")
	version := flag.Bool("v", false, "Показать версию и выйти")
	versionLong := flag.Bool("version", false, "Показать версию и выйти")
	
	flag.Parse()

	if *help || *helpLong {
		printHelp()
		os.Exit(0)
	}
	
	if *version || *versionLong {
		printVersion()
		os.Exit(0)
	}

	shellPid := os.Getppid() 
	syscall.Kill(shellPid, syscall.SIGHUP) 
}

func printHelp() {
	fmt.Println("Использование: exit [КЛЮЧ]")
	fmt.Println("Завершает текущий сеанс терминала.")
	fmt.Println()
	fmt.Println("Ключи:")
	fmt.Println("  -h, --help     Показать эту помощь и выйти")
	fmt.Println("  -v, --version  Показать версию и выйти")
	fmt.Println()
	fmt.Println("Пример:")
	fmt.Println("  exit     # Завершить текущий терминал")
	fmt.Println("  exit -h  # Показать справку")
}

func printVersion() {
	fmt.Println("exit версия 1.0")
	fmt.Println("Учебный проект - аналог системных утилит Linux")
	fmt.Println("Язык программирования: Golang")
}
