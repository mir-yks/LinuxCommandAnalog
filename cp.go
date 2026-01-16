package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Help        bool
	Version     bool
	Recursive   bool
	Interactive bool
	Verbose     bool
	Preserve    bool
	Update      bool
	Src         []string
	Dst         string
}

const version = "1.0.0"

func main() {
	config := parseArgs()
	
	if config.Help {
		printHelp()
		return
	}

	if config.Version {
		fmt.Println("cp версия 1.0")
		fmt.Println("Разработано в рамках учебного проекта")
		fmt.Println("Язык программирования: Golang")
		return
	}

	if len(config.Src) == 0 || config.Dst == "" {
		fmt.Fprintln(os.Stderr, "cp: укажите источник и назначение")
		printHelp()
		os.Exit(1)
	}

	if err := copyPaths(config); err != nil {
		fmt.Fprintf(os.Stderr, "cp: %v\n", err)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Println("Копирование завершено успешно.")
	}
}

func parseArgs() *Config {
	config := &Config{}
	
	flag.BoolVar(&config.Help, "h", false, "показать справку и выйти")
	flag.BoolVar(&config.Help, "help", false, "показать справку и выйти")
	flag.BoolVar(&config.Version, "version", false, "показать версию и выйти")
	flag.BoolVar(&config.Recursive, "r", false, "копировать рекурсивно (для каталогов)")
	flag.BoolVar(&config.Recursive, "R", false, "копировать рекурсивно (для каталогов)")
	flag.BoolVar(&config.Interactive, "i", false, "интерактивный режим (запрашивать перед перезаписью)")
	flag.BoolVar(&config.Verbose, "v", false, "подробный режим")
	flag.BoolVar(&config.Preserve, "p", false, "сохранять время изменения и права доступа")
	flag.BoolVar(&config.Update, "u", false, "копировать только новые файлы")
	
	flag.Parse()
	
	// Правильная обработка позиционных аргументов
	args := flag.Args()
	if len(args) == 0 {
		return config
	}
	
	config.Dst = args[len(args)-1]
	config.Src = args[:len(args)-1]
	
	return config
}

func printHelp() {
	fmt.Println("cp - копирует файлы и каталоги")
	fmt.Println()
	fmt.Println("Использование: cp [ОПЦИЯ]... ИСТОЧНИК... НАЗНАЧЕНИЕ")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -r, -R     рекурсивное копирование каталогов")
	fmt.Println("  -i         запрашивать подтверждение перед перезаписью")
	fmt.Println("  -v         подробный режим")
	fmt.Println("  -p         сохранять время изменения и права доступа")
	fmt.Println("  -u         копировать только новые файлы")
	fmt.Println("  -h, --help показать эту справку")
	fmt.Println("  --version  показать версию")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  cp file.txt backup.txt")
	fmt.Println("  cp -r dir1/ dir2/")
	fmt.Println("  cp *.txt /backup/")
}

func copyPaths(config *Config) error {
	for _, src := range config.Src {
		dstPath := config.Dst
		dstInfo, err := os.Stat(config.Dst)
		if err == nil && dstInfo.IsDir() {
			dstPath = filepath.Join(config.Dst, filepath.Base(src))
		}
		
		if err := copyPath(src, dstPath, config); err != nil {
			return err
		}
	}
	return nil
}

func copyPath(src, dst string, config *Config) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("источник '%s': %v", src, err)
	}

	if srcInfo.IsDir() {
		if !config.Recursive {
			return fmt.Errorf("омиттинг каталог '%s' (не указан -r/-R)", src)
		}
		return copyDir(src, dst, config)
	}

	return copyFile(src, dst, config)
}

func copyFile(src, dst string, config *Config) error {
	if config.Verbose {
		fmt.Printf("'%s' -> '%s'\n", src, dst)
	}

	if config.Update {
		srcInfo, _ := os.Stat(src)
		dstInfo, err := os.Stat(dst)
		if err == nil && !srcInfo.ModTime().After(dstInfo.ModTime()) {
			if config.Verbose {
				fmt.Printf("не копируется (устаревший): '%s'\n", src)
			}
			return nil
		}
	}

	if config.Interactive {
		if _, err := os.Stat(dst); err == nil {
			var response string
			fmt.Printf("%s уже существует. Перезаписать? (y/n): ", dst)
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				return nil
			}
		}
	}

	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать каталог %s: %v", dstDir, err)
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

	if config.Preserve {
		srcInfo, _ := os.Stat(src)
		os.Chmod(dst, srcInfo.Mode())
		os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime())
	}

	if config.Verbose {
		size := getFileSize(src)
		fmt.Printf("скопировано: %d байт\n", size)
	}

	return nil
}

func copyDir(srcDir, dstDir string, config *Config) error {
	if config.Verbose {
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

		return copyFile(srcPath, dstPath, config)
	})
}

func getFileSize(filename string) int64 {
	info, err := os.Stat(filename)
	if err != nil {
		return 0
	}
	return info.Size()
}

