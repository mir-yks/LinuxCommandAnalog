package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Help        bool
	Recursive   bool   
	Interactive bool   
	Verbose     bool   
	Src         string 
	Dst         string 
}

const version = "1.0.0"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "cp: критическая ошибка: %v\n", r)
			os.Exit(1)
		}
	}()

	config := parseArgs()
	
	if config.Help {
		printHelp()
		return
	}

	if config.Src == "" || config.Dst == "" {
		fmt.Fprintln(os.Stderr, "cp: укажите источник и назначение")
		printHelp()
		os.Exit(1)
	}

	if err := copyPath(config); err != nil {
		fmt.Fprintf(os.Stderr, "cp: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Копирование завершено успешно.")
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
			return config
		case "-r":
			config.Recursive = true
			i++
			continue
		case "-i":
			config.Interactive = true
			i++
			continue
		case "-v":
			config.Verbose = true
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
						return config
					case 'r':
						config.Recursive = true
					case 'i':
						config.Interactive = true
					case 'v':
						config.Verbose = true
					default:
						panic(fmt.Sprintf("неверная опция — '%s'", arg))
					}
				}
				i++
				continue
			}

			if config.Src == "" {
				config.Src = arg
				i++
				continue
			}
			if config.Dst == "" {
				config.Dst = arg
				i++
				continue
			}

			panic(fmt.Sprintf("неверный аргумент '%s'", arg))
		}
	}

	return config
}

func printHelp() {
	fmt.Println("cp - копирует файлы и каталоги")
	fmt.Println()
	fmt.Println("Использование: cp [ОПЦИЯ]... ИСТОЧНИК НАЗНАЧЕНИЕ")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -r     рекурсивное копирование каталогов")
	fmt.Println("  -i     запрашивать подтверждение перед перезаписью")
	fmt.Println("  -v     подробный режим")
	fmt.Println("  -h     показать эту справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  cp file.txt backup.txt")
	fmt.Println("  cp -r dir1/ dir2/")
	fmt.Println("  cp 1.txt var/2.txt")
}

func copyPath(config *Config) error {
	srcInfo, err := os.Stat(config.Src)
	if err != nil {
		return fmt.Errorf("источник недоступен: %v", err)
	}

	if !srcInfo.IsDir() {
		return copyFile(config.Src, config.Dst, config.Interactive, config.Verbose)
	}

	if !config.Recursive {
		return fmt.Errorf("копирование каталога требует -r")
	}

	dstInfo, err := os.Stat(config.Dst)
	if err == nil && dstInfo.IsDir() {
		return copyDir(config.Src, filepath.Join(config.Dst, filepath.Base(config.Src)), 
			config.Interactive, config.Verbose)
	}

	return copyDir(config.Src, config.Dst, config.Interactive, config.Verbose)
}

func copyFile(src, dst string, interactive, verbose bool) error {
	if verbose {
		fmt.Printf("'%s' -> '%s'\n", src, dst)
	}

	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать каталог %s: %v", dstDir, err)
	}

	if interactive {
		if _, err := os.Stat(dst); err == nil {
			var response string
			fmt.Printf("%s уже существует. Перезаписать? (y/n): ", dst)
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				return nil
			}
		}
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("не удалось открыть источник %s: %v", src, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %v", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("ошибка копирования данных: %v", err)
	}

	if verbose {
		size := getFileSize(src)
		fmt.Printf("скопировано: %d байт\n", size)
	}

	return nil
}


func copyDir(srcDir, dstDir string, interactive, verbose bool) error {
	if verbose {
		fmt.Printf("каталог '%s' -> '%s'\n", srcDir, dstDir)
	}

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать каталог %s: %v", dstDir, err)
	}

	return filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		if interactive {
			if _, err := os.Stat(dstPath); err == nil {
				var response string
				fmt.Printf("%s существует. Перезаписать? (y/n): ", dstPath)
				fmt.Scanln(&response)
				if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
					return nil
				}
			}
		}

		return copyFile(srcPath, dstPath, false, verbose)
	})
}

func getFileSize(filename string) int64 {
	info, err := os.Stat(filename)
	if err != nil {
		return 0
	}
	return info.Size()
}

