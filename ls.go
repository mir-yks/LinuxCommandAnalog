package main

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
	"time"
)

type Config struct {
	Help      bool   
	Version   bool   
	All       bool   
	Long      bool   
	Reverse   bool   
	Recursive bool   
	Paths     []string 
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

	if len(config.Paths) == 0 {
		config.Paths = append(config.Paths, ".")
	}

	executeLs(config) 
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
		case "-a":
			config.All = true
			i++
			continue
		case "-l":
			config.Long = true
			i++
			continue
		case "-r":
			config.Reverse = true
			i++
			continue
		case "-R":
			config.Recursive = true
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'a':
						config.All = true
					case 'l':
						config.Long = true
					case 'r':
						config.Reverse = true
					case 'R':
						config.Recursive = true
					case 'h':
						config.Help = true
					case 'v':
						config.Version = true
					default:
						panic(fmt.Sprintf("ls: неверный ключ — '%s'", arg))
					}
				}
				i++
				continue
			}
			config.Paths = append(config.Paths, arg) 
			i++
		}
	}

	return config
}

// printHelp выводит справку по использованию программы
func printHelp() {
	fmt.Println("ls - перечисление содержимого директорий")
	fmt.Println()
	fmt.Println("Использование: ls [ОПЦИЯ]... [ПУТЬ]...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -a     включить скрытые файлы")
	fmt.Println("  -l     длинный формат")
	fmt.Println("  -r     обратить порядок")
	fmt.Println("  -R     рекурсивный вывод")
	fmt.Println("  -h     показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  ls                     # Текущая директория")
	fmt.Println("  ls -l /tmp            # Длинный формат")
	fmt.Println("  ls -a -r dir          # Скрытые файлы, обратный порядок")
}

// printVersion выводит информацию о версии программы
func printVersion() {
	fmt.Println("ls версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeLs выполняет перечисление файлов для всех указанных путей
func executeLs(config *Config) {
	for _, path := range config.Paths {
		if err := listDirectory(path, config.All, config.Long, config.Reverse, config.Recursive); err != nil {
			fmt.Fprintf(os.Stderr, "ls: %v\n", err)
			os.Exit(1)
		}
	}
}

// FileDetails содержит всю информацию о файле для отображения
type FileDetails struct {
	Name    string      
	Size    int64       
	ModTime time.Time   
	Mode    fs.FileMode 
	IsDir   bool        
}

// listDirectory читает содержимое директории и подготавливает данные для вывода
func listDirectory(dirPath string, showHidden bool, longFormat bool, reverseOrder bool, recursive bool) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать '%s': %v", dirPath, err)
	}

	var files []FileDetails
	for _, entry := range entries {
		name := entry.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileDetails{
			Name:    name,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Mode:    info.Mode(),
			IsDir:   info.IsDir(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		if reverseOrder {
			return files[j].Name < files[i].Name
		}
		return files[i].Name < files[j].Name
	})

	if recursive {
		fmt.Printf("%s:\n", dirPath) 
	}

	displayFileInfo(files, showHidden, longFormat, dirPath, recursive)

	if recursive {
		for _, file := range files {
			if file.IsDir {
				fmt.Println() 
				if err := listDirectory(dirPath+"/"+file.Name, showHidden, longFormat, reverseOrder, recursive); err != nil {
					fmt.Fprintf(os.Stderr, "ls: %v\n", err)
				}
			}
		}
	}

	return nil
}

// displayFileInfo выводит файлы в нужном формате (одностолбцовый или многостолбцовый)
func displayFileInfo(files []FileDetails, showHidden bool, longFormat bool, dirPath string, recursive bool) {
	if longFormat {
		for _, file := range files {
			fmt.Printf("%s %6d %s %s\n",
				file.Mode,
				file.Size,
				file.ModTime.Format("Jan 02 15:04"),
				file.Name)
		}
		return
	}

	var visibleFiles []string
	for _, file := range files {
		if !showHidden && strings.HasPrefix(file.Name, ".") {
			continue
		}
		visibleFiles = append(visibleFiles, file.Name)
	}

	if len(visibleFiles) == 0 {
		return
	}

	const termWidth = 80 
	maxLen := 14
	for _, name := range visibleFiles {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	colWidth := maxLen + 2 
	cols := termWidth / colWidth
	if cols < 1 {
		cols = 1
	}

	rows := (len(visibleFiles) + cols - 1) / cols 

	grid := make([][]string, rows)
	for i := range grid {
		grid[i] = make([]string, cols)
	}

	for i, name := range visibleFiles {
		row := i / cols
		col := i % cols
		grid[row][col] = name
	}

	for _, row := range grid {
		line := ""
		for _, name := range row {
			if name != "" {
				padding := colWidth - len(name)
				line += name + strings.Repeat(" ", padding)
			}
		}
		if line != "" {
			fmt.Println(line[:len(line)-1]) 
		}
	}
}

