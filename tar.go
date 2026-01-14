package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	Help      bool     
	Version   bool     
	Create    bool     
	Extract   bool     
	Archive   string   
	Files     []string 
}

const ver = "1.0.0"

// printHelp выводит справку по использованию утилиты tar
func printHelp() {
	fmt.Println("tar - архивация файлов (tar.gz)")
	fmt.Println()
	fmt.Println("Использование: tar [ОПЦИЯ]... [ФАЙЛЫ]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -h, --help     Показать эту справку")
	fmt.Println("  -v, --version  Показать информацию о версии")
	fmt.Println("  -c             Создать архив")
	fmt.Println("  -x             Распаковать архив")
	fmt.Println("  -f ФАЙЛ       Имя архива (по умолчанию: archive.tar.gz)")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  tar -c -f archive.tar.gz file1.txt file2.txt")
	fmt.Println("  tar -x -f archive.tar.gz")
}

// printVersion выводит информацию о версии программы
func printVersion() {
	fmt.Println("tar версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// parseArgs разбирает аргументы командной строки
func parseArgs() *Config {
	config := &Config{Archive: "archive.tar.gz"}
	i := 1

	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-h", "--help":
			config.Help = true
			i++
			return config 
		case "-v", "--version":
			config.Version = true
			i++
			return config 
		case "-c":
			config.Create = true
			i++
			continue
		case "-x":
			config.Extract = true
			i++
			continue
		case "-f":
			if i+1 >= len(os.Args) {
				panic("tar: опция -f требует аргумент")
			}
			i++
			config.Archive = os.Args[i]
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				panic(fmt.Sprintf("tar: неверный ключ '%s'. Используйте --help", arg))
			}
			config.Files = append(config.Files, arg)
			i++
		}
	}

	if config.Create && config.Extract {
		panic("tar: не удается одновременно создавать и распаковывать")
	}
	if !config.Create && !config.Extract {
		panic("tar: не указано действие (-c или -x)")
	}

	return config
}

// createTarGz создает tar.gz архив
func createTarGz(archive string, files []string) error {
	file, err := os.Create(archive)
	if err != nil {
		return fmt.Errorf("не удается создать '%s': %v", archive, err)
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, filePath := range files {
		if err := addFileToTar(tw, filePath); err != nil {
			return err
		}
	}

	return nil
}

// addFileToTar рекурсивно добавляет файл/директорию в tar
func addFileToTar(tw *tar.Writer, filePath string) error {
	info, err := os.Lstat(filePath)
	if err != nil {
		return fmt.Errorf("не удается прочитать '%s': %v", filePath, err)
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filePath)

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if !info.IsDir() {
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractTarGz распаковывает tar.gz архив
func extractTarGz(archive string) error {
	file, err := os.Open(archive)
	if err != nil {
		return fmt.Errorf("не удается открыть '%s': %v", archive, err)
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("неверный gzip архив: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка чтения tar: %v", err)
		}

		if err := extractFile(tr, header); err != nil {
			return err
		}
	}

	return nil
}

// extractFile извлекает файл из tar
func extractFile(tr *tar.Reader, header *tar.Header) error {
	path := header.Name
	info := header.FileInfo()

	if info.IsDir() {
		return os.MkdirAll(path, info.Mode())
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, tr)
	return err
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "tar: %v\n", r)
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

	if config.Create {
		if len(config.Files) == 0 {
			panic("tar: не указаны файлы для архивации")
		}
		if err := createTarGz(config.Archive, config.Files); err != nil {
			panic(err)
		}
		fmt.Printf("Архив создан: %s\n", config.Archive)
	} else {
		if err := extractTarGz(config.Archive); err != nil {
			panic(err)
		}
		fmt.Printf("Архив распакован: %s\n", config.Archive)
	}
}

