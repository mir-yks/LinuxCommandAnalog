package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

const (
	letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
	symbols   = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	version   = "1.0.0"
)

type Config struct {
	Help           bool
	IncludeNumbers bool 
	IncludeSymbols  bool 
	AddSpecial     bool 
	Length         int  
	Count          int  
	NumColumns     int  
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "pwgen: критическая ошибка: %v\n", r)
			os.Exit(1)
		}
	}()

	config := parseArgs()
	
	if config.Help {
		printHelp()
		return
	}

	if config.Length <= 0 {
		config.Length = 8
	}
	if config.Count <= 0 {
		config.Count = 160
	}
	if config.NumColumns <= 0 {
		config.NumColumns = 8
	}

	passwords, err := generatePasswords(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwgen: ошибка генерации: %v\n", err)
		os.Exit(1)
	}

	displayPasswords(passwords, config)
}

// parseArgs разбирает аргументы командной строки
func parseArgs() *Config {
	config := &Config{Length: 8, Count: 160, NumColumns: 8}
	i := 1

	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-h":
			config.Help = true
			i++
			return config
		case "-n":
			config.IncludeNumbers = true
			i++
			continue
		case "-s":
			config.IncludeSymbols = true
			i++
			continue
		case "-y":
			config.AddSpecial = true
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
						return config
					case 'n':
						config.IncludeNumbers = true
					case 's':
						config.IncludeSymbols = true
					case 'y':
						config.AddSpecial = true
					default:
						panic(fmt.Sprintf("pwgen: неверная опция — '%s'", arg))
					}
				}
				i++
				continue
			}

			if length, err := strconv.Atoi(arg); err == nil && length > 0 {
				config.Length = length
				i++
				continue
			}

			panic(fmt.Sprintf("pwgen: неверный аргумент '%s'", arg))
		}
	}

	return config
}

func printHelp() {
	fmt.Println("pwgen - генератор безопасных паролей")
	fmt.Println()
	fmt.Println("Использование: pwgen [ОПЦИЯ]... [ДЛИНА] [КОЛИЧЕСТВО]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -n     включить цифры в пароль")
	fmt.Println("  -s     включить специальные символы") 
	fmt.Println("  -y     добавить один специальный символ в конец")
	fmt.Println("  -h     показать эту справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  pwgen                    # 160 паролей по 8 символов")
	fmt.Println("  pwgen 12                 # 160 паролей по 12 символов")
	fmt.Println("  pwgen 10 50              # 50 паролей по 10 символов")
	fmt.Println("  pwgen -n -s 16 100       # 100 паролей с цифрами/символами")
}

func generatePasswords(config *Config) ([]string, error) {
	passwords := make([]string, config.Count)
	charset := letters
	
	if config.IncludeNumbers {
		charset += numbers
	}
	if config.IncludeSymbols {
		charset += symbols
	}

	if len(charset) == 0 {
		return nil, fmt.Errorf("не задан набор символов")
	}

	for i := 0; i < config.Count; i++ {
		password, err := generateSinglePassword(charset, config.Length, config.AddSpecial)
		if err != nil {
			return nil, err
		}
		passwords[i] = password
	}

	return passwords, nil
}

func generateSinglePassword(charset string, length int, addSpecial bool) (string, error) {
	password := make([]byte, length)
	
	for i := range password {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("ошибка генерации: %v", err)
		}
		password[i] = charset[index.Int64()]
	}

	if addSpecial && len(symbols) > 0 {
		specialIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
		if err != nil {
			return "", fmt.Errorf("ошибка спецсимвола: %v", err)
		}
		password = append(password, symbols[specialIndex.Int64()])
	}

	return string(password), nil
}

func displayPasswords(passwords []string, config *Config) {
	const passwordsPerLine = 8
	padding := config.Length + 1
	
	for i := 0; i < len(passwords); i += passwordsPerLine {
		line := ""
		end := i + passwordsPerLine
		if end > len(passwords) {
			end = len(passwords)
		}
		
		for j := i; j < end; j++ {
			padded := fmt.Sprintf("%-*s", padding, passwords[j])
			line += padded
		}
		fmt.Println(line)
	}
}

