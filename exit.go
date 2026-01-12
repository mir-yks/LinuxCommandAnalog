package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

func main() {
	// Обработка ключа -help
	help := flag.Bool("help", false, "Показать помощь и выйти")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	shellPid :=os.Getppid() // получение PID терминала
	syscall.Kill(shellPid, syscall.SIGHUP) 
}

func printHelp() {
	fmt.Println("Использование: exit")
	fmt.Println("Ключи:")
	fmt.Println("  -help       Показать это сообщение и выйти")
}
