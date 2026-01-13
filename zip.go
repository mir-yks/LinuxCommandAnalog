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
	Version   bool
	Delete    bool  
	Update    bool  
	Archive   string
	Filenames []string
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

	if config.Archive == "" || len(config.Filenames) == 0 {
		fmt.Fprintln(os.Stderr, "zip: пропущен архив или файлы")
		fmt.Fprintln(os.Stderr, "По команде «zip -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	executeZip(config)
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
		case "-d":
			config.Delete = true
			i++
			continue
		case "-u":
			config.Update = true
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'v':
						config.Version = true
					case 'd':
						config.Delete = true
					case 'u':
						config.Update = true
					default:
						panic(fmt.Sprintf("zip: неверный ключ — '%s'", arg))
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

			config.Filenames = append(config.Filenames, arg)
			i++
		}
	}

	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("zip - создает ZIP архивы")
	fmt.Println()
	fmt.Println("Использование: zip [ОПЦИЯ]... АРХИВ ФАЙЛЫ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -d     удалить файлы после архивации")
	fmt.Println("  -u     обновить существующий архив (добавить файлы)")
	fmt.Println("  -h     показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  zip archive.zip file1.txt file2.txt")
	fmt.Println("  zip -u archive.zip new.txt         # Добавить в существующий архив")
	fmt.Println("  zip -d archive.zip file.txt        # Архивировать и удалить")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("zip версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeZip выполняет архивацию
func executeZip(config *Config) {
	var finalErr error
	
	if config.Update {
		finalErr = updateArchiveViaTemp(config)
	} else {
		zipWriter, file, err := createZipWriter(config.Archive)
		if err != nil {
			fmt.Fprintf(os.Stderr, "zip: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		for _, filename := range config.Filenames {
			if err := addFileToZip(zipWriter, filename, "", false); err != nil {
				fmt.Fprintf(os.Stderr, "zip: %v\n", err)
			}
			if config.Delete {
				os.Remove(filename)
			}
		}

		if err := zipWriter.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "zip: ошибка завершения архива: %v\n", err)
			os.Exit(1)
		}
	}

	if finalErr != nil {
		fmt.Fprintf(os.Stderr, "zip: %v\n", finalErr)
		os.Exit(1)
	}
}

// createZipWriter создает новый ZIP архив
func createZipWriter(archiveName string) (*zip.Writer, *os.File, error) {
	file, err := os.Create(archiveName)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось создать архив: %v", err)
	}
	return zip.NewWriter(file), file, nil
}

// updateArchiveViaTemp обновляет архив через временный файл
func updateArchiveViaTemp(config *Config) error {
	tempArchive := config.Archive + ".tmp"
	
	zipWriter, tempFile, err := createZipWriter(tempArchive)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	oldZipFile, err := os.Open(config.Archive)
	if err != nil {
		return fmt.Errorf("не удалось открыть старый архив: %v", err)
	}
	defer oldZipFile.Close()

	fileInfo, err := oldZipFile.Stat()
	if err != nil {
		return fmt.Errorf("не удалось получить размер архива: %v", err)
	}

	oldZipReader, err := zip.NewReader(oldZipFile, fileInfo.Size())
	if err != nil {
		return fmt.Errorf("не удалось прочитать старый архив: %v", err)
	}

	for _, fileHeader := range oldZipReader.File {
		writer, err := zipWriter.Create(fileHeader.Name)
		if err != nil {
			return fmt.Errorf("не удалось скопировать '%s': %v", fileHeader.Name, err)
		}

		oldFile, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("не удалось открыть '%s': %v", fileHeader.Name, err)
		}
		defer oldFile.Close()

		_, err = io.Copy(writer, oldFile)
		if err != nil {
			return fmt.Errorf("не удалось скопировать '%s': %v", fileHeader.Name, err)
		}
	}

	for _, filename := range config.Filenames {
		if err := addFileToZip(zipWriter, filename, "", false); err != nil {
			return err
		}
		if config.Delete {
			os.Remove(filename)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("ошибка завершения временного архива: %v", err)
	}

	if err := os.Rename(tempArchive, config.Archive); err != nil {
		os.Remove(tempArchive)
		return fmt.Errorf("не удалось заменить архив: %v", err)
	}

	return nil
}

// addFileToZip добавляет файл в архив
func addFileToZip(zipWriter *zip.Writer, filePath string, basePath string, verbose bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть '%s': %v", filePath, err)
	}
	defer file.Close()

	zipEntry := filepath.Base(filePath)
	if basePath != "" {
		zipEntry = filepath.Join(filepath.Base(basePath), filepath.Base(filePath))
	}

	writer, err := zipWriter.Create(zipEntry)
	if err != nil {
		return fmt.Errorf("не удалось создать запись '%s': %v", zipEntry, err)
	}

	if verbose {
		fmt.Printf("Добавлен: %s\n", zipEntry)
	}

	_, err = io.Copy(writer, file)
	return err
}

