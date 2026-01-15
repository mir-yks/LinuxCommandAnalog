package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Help     bool
	All      bool   
	User     bool   
	Group    bool   
	Current  bool   
	Username string 
}

const version = "1.0.0"

type ProcInfo struct {
	User    string
	PID     int
	CPU     float64
	Mem     float64
	VSZ     int
	RSS     int
	TTY     string
	STAT    string
	Start   string
	Time    string
	Command string
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "ps: критическая ошибка: %v\n", r)
			os.Exit(1)
		}
	}()

	config := parseArgs()
	
	if config.Help {
		printHelp()
		return
	}

	procList, err := fetchProcesses(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ps: %v\n", err)
		os.Exit(1)
	}

	displayProcesses(procList, config)
}

// parseArgs разбирает аргументы командной строки
func parseArgs() *Config {
	config := &Config{Current: true}
	i := 1

	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-h":
			config.Help = true
			i++
			return config
		case "-a":
			config.All = true
			config.Current = false
			i++
			continue
		case "-u":
			config.User = true
			config.Current = false
			i++
			continue
		case "-g":
			config.Group = true
			config.Current = false
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				for _, ch := range arg[1:] {
					switch ch {
					case 'h':
						config.Help = true
						return config
					case 'a':
						config.All = true
						config.Current = false
					case 'u':
						config.User = true
						config.Current = false
					case 'g':
						config.Group = true
						config.Current = false
					default:
						panic(fmt.Sprintf("неверная опция — '%s'", arg))
					}
				}
				i++
				continue
			}

			config.Username = arg
			config.Current = false
			i++
		}
	}

	return config
}

func printHelp() {
	fmt.Println("ps - отображает информацию о процессах")
	fmt.Println()
	fmt.Println("Использование: ps [ОПЦИЯ]... [ПОЛЬЗОВАТЕЛЬ]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -a     показать информацию о всех процессах")
	fmt.Println("  -u     показать процессы текущего пользователя") 
	fmt.Println("  -g     показать все процессы, включая фоновые")
	fmt.Println("  -h     показать эту справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  ps                    # процессы текущего пользователя")
	fmt.Println("  ps -a                 # все процессы")
	fmt.Println("  ps -u root            # процессы пользователя root")
	fmt.Println("  ps -g                 # все процессы с подробностями")
}

func fetchProcesses(config *Config) ([]ProcInfo, error) {
	var procList []ProcInfo

	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать /proc: %v", err)
	}

	currentUser, _ := user.Current()
	
	for _, file := range files {
		if !file.IsDir() || !isNumeric(file.Name()) {
			continue
		}

		pidStr := file.Name()
		
		procInfo, err := parseProcInfo(pidStr)
		if err != nil {
			continue
		}

		if shouldIncludeProcess(config, procInfo, currentUser) {
			procList = append(procList, procInfo)
		}
	}

	return procList, nil
}

func shouldIncludeProcess(config *Config, proc ProcInfo, currentUser *user.User) bool {
	if config.All {
		return true
	}

	if config.Current {
		return proc.User == currentUser.Username
	}

	if config.Username != "" {
		return proc.User == config.Username
	}

	return true
}

func parseProcInfo(pidStr string) (ProcInfo, error) {
	pidInt, err := strconv.Atoi(pidStr)
	if err != nil {
		return ProcInfo{}, fmt.Errorf("неверный PID: %s", pidStr)
	}

	statPath := filepath.Join("/proc", pidStr, "stat")
	statusPath := filepath.Join("/proc", pidStr, "status")
	cmdlinePath := filepath.Join("/proc", pidStr, "cmdline")

	statBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return ProcInfo{}, nil
	}
	
	fields := strings.Fields(string(statBytes))
	if len(fields) < 24 {
		return ProcInfo{}, nil
	}
	
	statusBytes, err := ioutil.ReadFile(statusPath)
	if err != nil {
		return ProcInfo{}, nil
	}
	userName := extractUserFromStatus(statusBytes)

	cmdlineBytes, err := ioutil.ReadFile(cmdlinePath)
	var command string
	if err == nil {
		command = strings.ReplaceAll(string(cmdlineBytes), "\x00", " ")
		command = strings.TrimSpace(command)
	}
	
	if command == "" {
		command = fields[1]
		if len(command) > 1 && command[0] == '(' && command[len(command)-1] == ')' {
			command = command[1 : len(command)-1]
		}
	}

	proc := ProcInfo{
		User:    userName,
		PID:     pidInt,  
		TTY:     fields[6],
		STAT:    fields[2],
		Time:    fields[13],
		Start:   fields[21],
		VSZ:     extractIntField(fields, 22),
		RSS:     extractIntField(fields, 23),
		CPU:     calculateCPUUsage(fields),
		Mem:     calculateMemoryUsage(fields),
		Command: command,
	}

	return proc, nil
}

func extractUserFromStatus(status []byte) string {
	lines := strings.Split(string(status), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Uid:") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				uidInt, err := strconv.Atoi(fields[1])
				if err == nil {
					userInfo, err := user.LookupId(strconv.Itoa(uidInt))
					if err == nil {
						return userInfo.Username
					}
				}
			}
		}
	}
	return "unknown"
}

func calculateCPUUsage(fields []string) float64 {
	if len(fields) > 13 {
		utime, _ := strconv.ParseFloat(fields[13], 64)
		return utime / 100.0 
	}
	return 0.0
}

func calculateMemoryUsage(fields []string) float64 {
	if len(fields) > 23 {
		rss := extractIntField(fields, 23)
		return float64(rss) * 4 / 1024.0 
	}
	return 0.0
}

func extractIntField(fields []string, index int) int {
	if index < len(fields) {
		if val, err := strconv.Atoi(fields[index]); err == nil {
			return val
		}
	}
	return 0
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func displayProcesses(procList []ProcInfo, config *Config) {
	if config.All {
		fmt.Println("PID TTY TIME CMD")
		for _, proc := range procList {
			fmt.Printf("%d %s %s %s\n", proc.PID, proc.TTY, proc.Time, proc.Command)
		}
	} else if config.Group {
		fmt.Println("PID TTY STAT TIME COMMAND")
		for _, proc := range procList {
			fmt.Printf("%d %s %s %s %s\n", proc.PID, proc.TTY, proc.STAT, proc.Time, proc.Command)
		}
	} else {
		fmt.Println("USER PID %CPU %MEM VSZ RSS TTY STAT START TIME COMMAND")
		for _, proc := range procList {
			fmt.Printf("%s %d %.1f %.1f %d %d %s %s %s %s %s\n",
				proc.User, proc.PID, proc.CPU, proc.Mem, proc.VSZ, proc.RSS,
				proc.TTY, proc.STAT, proc.Start, proc.Time, proc.Command)
		}
	}
}

