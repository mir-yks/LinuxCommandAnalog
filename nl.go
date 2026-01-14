package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Help      bool   
	BodyMode  string 
	NumberFmt string 
	Width     int    
	FilePath  string 
}

// printHelp выводит справку по использованию утилиты nl
func printHelp() {
	fmt.Println("nl - нумерация строк")
	fmt.Println()
	fmt.Println("Использование: nl [ОПЦИЯ]... [ФАЙЛ]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -h        Показать эту справку")
	fmt.Println("  -b <a|t>  Режим нумерации: a - все строки, t - только непустые")
	fmt.Println("  -n <ln|rn> Формат номера: ln - слева, rn - справа")
	fmt.Println("  -w <число> Ширина поля номера (1-20 символов)")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  nl file.txt              # Нумеровать все строки")
	fmt.Println("  nl -b t -w 4 file.txt    # Только непустые строки, ширина 4")
	fmt.Println("  nl -n rn file.txt        # Номера справа")
}

// parseArgs разбирает аргументы командной строки
func parseArgs() *Config {
	config := &Config{
		BodyMode: "a",  
		NumberFmt: "ln", 
		Width:     6,    
	}

	i := 1
	for i < len(os.Args) {
		arg := os.Args[i]
		switch arg {
		case "-h":
			config.Help = true
			return config
		case "-b":
			if i+1 >= len(os.Args) {
				panic("nl: опция -b требует аргумент")
			}
			i++
			mode := strings.TrimSpace(os.Args[i])
			if mode != "a" && mode != "t" {
				panic("nl: неверный режим -b. Используйте 'a' или 't'")
			}
			config.BodyMode = mode
			i++
		case "-n":
			if i+1 >= len(os.Args) {
				panic("nl: опция -n требует аргумент")
			}
			i++
			fmtArg := strings.TrimSpace(os.Args[i])
			if fmtArg != "ln" && fmtArg != "rn" {
				panic("nl: неверный формат -n. Используйте 'ln' или 'rn'")
			}
			config.NumberFmt = fmtArg
			i++
		case "-w":
			if i+1 >= len(os.Args) {
				panic("nl: опция -w требует аргумент")
			}
			i++
			widthStr := strings.TrimSpace(os.Args[i])
			width, err := strconv.Atoi(widthStr)
			if err != nil {
				panic("nl: ширина -w должна быть числом")
			}
			if width < 1 || width > 20 {
				panic("nl: ширина -w должна быть от 1 до 20")
			}
			config.Width = width
			i++
		default:
			if config.FilePath != "" {
				panic(fmt.Sprintf("nl: несколько файлов: %s", arg))
			}
			config.FilePath = arg
			i++
		}
	}

	return config
}

// numberLines нумерует строки согласно конфигурации
func numberLines(scanner *bufio.Scanner, config *Config) {
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		
		if config.BodyMode == "t" && strings.TrimSpace(line) == "" {
			fmt.Println(line)
			continue
		}

		var numberStr string
		if config.NumberFmt == "ln" {
			numberStr = fmt.Sprintf("%*d", config.Width, lineNumber)
		} else { 
			numberStr = fmt.Sprintf("%-*d", config.Width, lineNumber)
		}

		fmt.Printf("%s  %s\n", numberStr, line)
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("nl: ошибка чтения: %v", err))
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "nl: %v\n", r)
			os.Exit(1)
		}
	}()

	config := parseArgs()

	if config.Help {
		printHelp()
		return
	}

	var file *os.File
	var err error
	
	if config.FilePath != "" {
		file, err = os.Open(config.FilePath)
		if err != nil {
			panic(fmt.Sprintf("не удается открыть '%s': %v", config.FilePath, err))
		}
		defer file.Close()
	} else {
		file = os.Stdin
	}

	scanner := bufio.NewScanner(file)
	numberLines(scanner, config)
}
