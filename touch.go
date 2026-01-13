package main

import (
	"fmt"
	"os"
	"time"
)

const ver = "1.0.0"

type Config struct {
	Help      bool
	Version   bool
	Access    bool
	Modify    bool
	Filenames []string
}

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

	if len(config.Filenames) == 0 {
		fmt.Fprintln(os.Stderr, "touch: пропущен операнд, задающий файл")
		fmt.Fprintln(os.Stderr, "По команде «touch -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	err := executeTouch(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() *Config {
	config := &Config{}
	filenames := []string{}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-h":
			config.Help = true
		case "-v":
			config.Version = true
		case "--version":
			config.Version = true
		case "-a":
			config.Access = true
		case "-m":
			config.Modify = true
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'v':
						config.Version = true
					case 'a':
						config.Access = true
					case 'm':
						config.Modify = true
					default:
						panic(fmt.Sprintf("touch: неверный ключ — '%s'", arg))
					}
				}
			} else {
				filenames = append(filenames, arg)
			}
		}
	}

	config.Filenames = filenames
	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("touch - изменяет временные метки файлов или создает пустые файлы")
	fmt.Println()
	fmt.Println("Использование: touch [ОПЦИЯ]... ФАЙЛ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -a        Изменяет только время последнего доступа к файлу")
	fmt.Println("  -m        Изменяет только время последней модификации файла")
	fmt.Println("  -h        Показать эту справку")
	fmt.Println("  -v, --version Показать информацию о версии")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  touch file.txt           # Создает файл или обновляет временные метки")
	fmt.Println("  touch -a file.txt        # Обновляет только время доступа")
	fmt.Println("  touch -m file.txt        # Обновляет только время модификации")
	fmt.Println("  touch -am file.txt       # Обновляет обе временные метки (по умолчанию)")
	fmt.Println("  touch file1.txt file2.txt # Работа с несколькими файлами")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("touch версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeTouch выполняет основную логику команды touch
func executeTouch(config *Config) error {
	now := time.Now()

	for _, filename := range config.Filenames {
		err := processFile(filename, config, now)
		if err != nil {
			fmt.Fprintf(os.Stderr, "touch: ошибка при обработке файла '%s': %v\n", filename, err)
		}
	}

	return nil
}

// processFile обрабатывает один файл
func processFile(filename string, config *Config, now time.Time) error {
	fileInfo, err := os.Stat(filename)

	if os.IsNotExist(err) {
		file, createErr := os.Create(filename)
		if createErr != nil {
			return fmt.Errorf("не удалось создать файл: %v", createErr)
		}
		file.Close()
		
		fileInfo, err = os.Stat(filename)
		if err != nil {
			return fmt.Errorf("не удалось получить информацию о созданном файле: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("не удалось получить информацию о файле: %v", err)
	}

	// Определяем, какие временные метки обновлять
	atime := fileInfo.ModTime() 
	mtime := fileInfo.ModTime() 

	if config.Access && config.Modify {
		atime = now
		mtime = now
	} else if config.Access {
		atime = now
	} else if config.Modify {
		mtime = now
	} else {
		atime = now
		mtime = now
	}

	err = os.Chtimes(filename, atime, mtime)
	if err != nil {
		return fmt.Errorf("не удалось обновить временные метки: %v", err)
	}

	fmt.Printf("touch: обработан файл '%s'\n", filename)
	return nil
}

