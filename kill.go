package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

type Config struct {
	Help     bool
	Version  bool
	Pid      int
	Signal   string
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

	if config.Pid == 0 {
		fmt.Fprintln(os.Stderr, "kill: пропущен PID процесса")
		fmt.Fprintln(os.Stderr, "По команде «kill -h» можно получить дополнительную информацию.")
		os.Exit(1)
	}

	executeKill(config)
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
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
					case 'v':
						config.Version = true
					default:
						panic(fmt.Sprintf("kill: неверный ключ — '%s'", arg))
					}
				}
				i++
				continue
			}

			if config.Pid == 0 {
				pid, err := strconv.Atoi(arg)
				if err != nil || pid <= 0 {
					panic(fmt.Sprintf("kill: неверный PID '%s'", arg))
				}
				config.Pid = pid
				i++
				continue
			}

			if config.Signal == "" {
				config.Signal = arg
				i++
				continue
			}

			panic(fmt.Sprintf("kill: лишний аргумент '%s'", arg))
		}
	}

	return config
}

// printHelp выводит справку
func printHelp() {
	fmt.Println("kill - отправляет сигнал процессам")
	fmt.Println()
	fmt.Println("Использование: kill [ОПЦИЯ]... PID [СИГНАЛ]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -h        показать эту справку")
	fmt.Println("  -v, --version показать информацию о версии")
	fmt.Println()
	fmt.Println("Сигналы:")
	fmt.Println("  SIGTERM   завершение (по умолчанию)")
	fmt.Println("  SIGKILL   принудительное завершение")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  kill 1234                    # SIGTERM по умолчанию")
	fmt.Println("  kill 1234 SIGKILL            # Принудительное завершение")
	fmt.Println("  kill 1234 15                  # Сигнал 15 (SIGTERM)")
}

// printVersion выводит информацию о версии
func printVersion() {
	fmt.Println("kill версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// executeKill выполняет отправку сигнала процессу
func executeKill(config *Config) {
	signal := syscall.SIGTERM 
	numSignal := 15           

	if config.Signal != "" {
		if sigNum, err := strconv.Atoi(config.Signal); err == nil {
			signal = syscall.Signal(sigNum)
		} else {
			switch config.Signal {
			case "SIGTERM", "TERM", "15":
				signal = syscall.SIGTERM
			case "SIGKILL", "KILL", "9":
				signal = syscall.SIGKILL
			default:
				fmt.Fprintf(os.Stderr, "kill: неверный сигнал '%s'\n", config.Signal)
				os.Exit(1)
			}
		}
		numSignal = int(signal)
	}

	err := syscall.Kill(config.Pid, signal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kill: отправка сигнала %s процессу %d: %v\n", signal, config.Pid, err)
		os.Exit(1)
	}

	fmt.Printf("Сигнал %s (%d) отправлен процессу %d\n", signal, numSignal, config.Pid)
}
