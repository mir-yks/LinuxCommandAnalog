package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	Help      bool
	List      bool
	Archive   string
	OutputDir string
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

	if config.List {
		if config.Archive == "" {
			fmt.Fprintln(os.Stderr, "unzip: пропущен архив")
			fmt.Fprintln(os.Stderr, "По команде «unzip -h» можно получить дополнительную информацию.")
			os.Exit(1)
		}
		executeList(config)
		return
	}

	if config.Archive == "" {
		fmt.Fprintln(os.Stderr, "unzip: пропущен архив")
		fmt.Fprintln(os.Stderr, "По команде «unzip -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	executeExtract(config)
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
		case "-l":
			config.List = true
			i++
			continue
		case "-f":
			if i+1 < len(os.Args) {
				config.Archive = os.Args[i+1]
				i += 2
				continue
			} else {
				panic("не указан архив после -f")
			}
		case "-o":
			if i+1 < len(os.Args) {
				config.OutputDir = os.Args[i+1]
				i += 2
				continue
			} else {
				panic("не указана папка после -o")
			}
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'l':
						config.List = true
					case 'f':
						panic("параметр -f требует указания архива")
					case 'o':
						panic("параметр -o требует указания папки")
					default:
						panic(fmt.Sprintf("unzip: неверный ключ — '%s'", arg))
					}
				}
				i++
				continue
			}
			if config.Archive == "" {
				config.Archive = arg
				i++
				continue
			}
			i++
		}
	}

	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("unzip - извлекает файлы из ZIP архивов")
	fmt.Println()
	fmt.Println("Использование: unzip [ОПЦИЯ]... [АРХИВ] [-o ПАПКА]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  АРХИВ     zip-архив для извлечения (обязательно)")
	fmt.Println("  -o ПАПКА  папка для извлечения файлов (по умолчанию: текущая папка)")
	fmt.Println("  -l        вывести список файлов в архиве")
	fmt.Println("  -f        явно указать архив (альтернатива позиционному аргументу)")
	fmt.Println("  -h        показать эту справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  unzip archive.zip                    # В текущую папку")
	fmt.Println("  unzip -f archive.zip -o ./extracted  # В указанную папку")
	fmt.Println("  unzip -l archive.zip                 # Показать содержимое")
}


// executeList выводит список файлов в архиве
func executeList(config *Config) {
	reader, err := zip.OpenReader(config.Archive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unzip: %v\n", err)
		os.Exit(1)
	}
	defer reader.Close()

	fmt.Println("Список файлов в архиве:")
	for _, file := range reader.File {
		fmt.Println(file.Name)
	}
}

// executeExtract извлекает файлы из архива
func executeExtract(config *Config) {
	if err := unzipArchive(config.Archive, config.OutputDir); err != nil {
		fmt.Fprintf(os.Stderr, "unzip: %v\n", err)
		os.Exit(1)
	}
	
	outputFolder := config.OutputDir
	if config.OutputDir == "." {
		cwd, err := os.Getwd()
		if err == nil {
			outputFolder = filepath.Base(cwd)
		}
	}
	
	fmt.Printf("Файлы успешно извлечены в папку: %s\n", outputFolder)
}


// unzipArchive извлекает файлы из ZIP-архива в указанную папку
func unzipArchive(zipFile, destination string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("ошибка при открытии архива: %v", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		outputPath := filepath.Join(destination, file.Name)

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("ошибка при создании директорий для '%s': %v", outputPath, err)
		}

		if file.FileInfo().IsDir() {
			continue
		}

		inputFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("ошибка при открытии '%s' в архиве: %v", file.Name, err)
		}
		defer inputFile.Close()

		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("ошибка при создании '%s': %v", outputPath, err)
		}
		defer outputFile.Close()

		if _, err := io.Copy(outputFile, inputFile); err != nil {
			return fmt.Errorf("ошибка при копировании '%s': %v", file.Name, err)
		}
	}

	return nil
}

