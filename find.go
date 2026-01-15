package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Help        bool
	Directory   string 
	NamePattern string 
	MinSize     int64  
}

const version = "1.0.0"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "find: критическая ошибка: %v\n", r)
			os.Exit(1)
		}
	}()

	config := parseArgs()
	
	if config.Help {
		printHelp()
		return
	}

	if config.Directory == "" {
		config.Directory = "."
	}

	if config.NamePattern == "" {
		config.NamePattern = "*" 
	}

	if err := findFiles(config); err != nil {
		fmt.Fprintf(os.Stderr, "find: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() *Config {
	config := &Config{MinSize: 0}
	i := 1

	if i < len(os.Args) && !strings.HasPrefix(os.Args[i], "-") {
		config.Directory = os.Args[i]
		i++
	}

	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-h":
			config.Help = true
			i++
			return config
		case "-d":
			i++
			if i >= len(os.Args) {
				panic("ожидается путь после -d")
			}
			config.Directory = os.Args[i]
			i++
			continue
		case "-n":
			i++
			if i >= len(os.Args) {
				panic("ожидается шаблон после -n")
			}
			config.NamePattern = os.Args[i]
			i++
			continue
		case "-s":
			i++
			if i >= len(os.Args) {
				panic("ожидается размер после -s")
			}
			size, err := strconv.ParseInt(os.Args[i], 10, 64)
			if err != nil || size < 0 {
				panic(fmt.Sprintf("неверный размер: %s", os.Args[i]))
			}
			config.MinSize = size
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
						return config
					case 'd', 'n', 's':
						continue
					default:
						panic(fmt.Sprintf("неверная опция — '%s'", arg))
					}
				}
				i++
				continue
			}

			if config.NamePattern == "" {
				config.NamePattern = arg
				i++
				continue
			}

			panic(fmt.Sprintf("неверный аргумент '%s'", arg))
		}
	}

	return config
}

func printHelp() {
	fmt.Println("find - ищет файлы по имени и размеру")
	fmt.Println()
	fmt.Println("Использование: find [ДИРЕКТОРИЯ] [ШАБЛОН] [ОПЦИИ]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -d ДИРЕКТОРИЯ   директория поиска (по умолчанию \".\")")
	fmt.Println("  -n ШАБЛОН       имя или шаблон (*.go, test*, *log)")
	fmt.Println("  -s РАЗМЕР       минимальный размер в байтах")
	fmt.Println("  -h              показать справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  find -n '*.go'              # все .go файлы")
	fmt.Println("  find -d ./tmp -n '*.log'     # .log в /tmp")
	fmt.Println("  find -s 1024                # файлы > 1KB")
	fmt.Println("  find -h                     # справка")
}

func findFiles(config *Config) error {
	return filepath.WalkDir(config.Directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка доступа %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка информации %s: %v\n", path, err)
			return nil
		}

		if fileInfo.Size() < config.MinSize {
			return nil
		}

		if matchesPattern(fileInfo.Name(), config.NamePattern) {
			fmt.Println(path)
		}

		return nil
	})
}

func matchesPattern(filename, pattern string) bool {
	pattern = strings.TrimSpace(pattern)
	
	if pattern == "" || pattern == "*" {
		return true
	}
	
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(filename, prefix)
	}
	
	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(filename, suffix)
	}
	
	if strings.HasPrefix(pattern, "*.") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(filename, suffix)
	}
	
	return filename == pattern
}

