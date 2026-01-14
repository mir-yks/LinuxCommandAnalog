package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Config struct {
	Help     bool   
	Version  bool   
	All      bool   
	BlockSize string 
	Direct   bool   
}

const ver = "1.0.0"

type BlockInfo struct {
	Size   uint64 
	Suffix string 
}

// printHelp выводит справку по использованию утилиты df
func printHelp() {
	fmt.Println("df - отчет о использовании дискового пространства")
	fmt.Println()
	fmt.Println("Использование: df [ОПЦИЯ]... [ПУТЬ]")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -h, --help     Показать эту справку")
	fmt.Println("  -v, --version  Показать информацию о версии")
	fmt.Println("  -a             Показать все файловые системы")
	fmt.Println("  -B <размер>    Указать размер (K,M,G)")
	fmt.Println("  --direct       Прямые данные без кэша")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  df                       # Информация о корневой ФС")
	fmt.Println("  df -a                   # Все файловые системы")
	fmt.Println("  df -B 1M                # Размер в мегабайтах")
}

// printVersion выводит информацию о версии программы
func printVersion() {
	fmt.Println("df версия", ver)
	fmt.Println("Разработано в рамках учебного проекта")
	fmt.Println("Язык программирования: Golang")
}

// parseArgs разбирает аргументы командной строки
func parseArgs() *Config {
	config := &Config{}
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
		case "-a":
			config.All = true
			i++
			continue
		case "-B":
			if i+1 >= len(os.Args) {
				panic("df: опция -B требует аргумент")
			}
			i++
			config.BlockSize = os.Args[i]
			i++
			continue
		case "--direct":
			config.Direct = true
			i++
			continue
		default:
			if len(arg) > 1 && arg[0] == '-' {
				panic(fmt.Sprintf("df: неверный ключ '%s'. Используйте --help для справки", arg))
			}
			i++
		}
	}

	return config
}

// parseSize преобразует строку размера в BlockInfo
func parseSize(sizeStr string) (BlockInfo, error) {
	sizeStr = strings.TrimSpace(strings.ToUpper(sizeStr))
	
	if sizeStr == "" {
		return BlockInfo{Size: 1024, Suffix: "K"}, nil 
	}

	numStr := sizeStr
	suffix := ""
	
	lastChar := sizeStr[len(sizeStr)-1]
	if lastChar == 'K' || lastChar == 'M' || lastChar == 'G' {
		numStr = sizeStr[:len(sizeStr)-1]
		suffix = string(lastChar)
	}

	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return BlockInfo{}, fmt.Errorf("неверный размер: %s", sizeStr)
	}

	blockSize := num
	switch suffix {
	case "K":
		blockSize *= 1024
	case "M":
		blockSize *= 1024 * 1024
	case "G":
		blockSize *= 1024 * 1024 * 1024
	}

	return BlockInfo{Size: blockSize, Suffix: suffix}, nil
}

// getFSInfo получает информацию о файловых системах
func getFSInfo(all, direct bool) ([]struct {
	MountPoint string
	Total      uint64
	Used       uint64
	Free       uint64
}, error) {
	var fsStats []struct {
		MountPoint string
		Total      uint64
		Used       uint64
		Free       uint64
	}

	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return nil, fmt.Errorf("не удается прочитать /proc/mounts: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		mountPoint := parts[1]
		var stat syscall.Statfs_t
		
		if err := syscall.Statfs(mountPoint, &stat); err != nil {
			continue
		}

		blockSize := uint64(stat.Bsize)
		total := stat.Blocks * blockSize
		free := stat.Bfree * blockSize
		used := total - free

		fsStats = append(fsStats, struct {
			MountPoint string
			Total      uint64
			Used       uint64
			Free       uint64
		}{
			MountPoint: mountPoint,
			Total:      total,
			Used:       used,
			Free:       free,
		})
	}

	return fsStats, nil
}

// executeDf выполняет основную логику утилиты df
func executeDf(config *Config) {
	blockInfo, err := parseSize(config.BlockSize)
	if err != nil {
		panic(fmt.Sprintf("df: %v", err))
	}

	stats, err := getFSInfo(config.All, config.Direct)
	if err != nil {
		panic(fmt.Sprintf("df: %v", err))
	}

	suffix := blockInfo.Suffix
	if suffix == "" {
		suffix = ""
	}
	fmt.Printf("%-30s %12s %12s %12s\n", 
		"Файловая система", 
		"Размер"+suffix, 
		"Использовано"+suffix, 
		"Доступно"+suffix)

	for _, stat := range stats {
		total := stat.Total / blockInfo.Size
		used := stat.Used / blockInfo.Size
		free := stat.Free / blockInfo.Size

		if config.All || total > 0 {
			fmt.Printf("%-30s %12d %12d %12d\n", stat.MountPoint, total, used, free)
		}
	}

	if config.Direct {
		fmt.Println("Прямые данные получены через syscall.Statfs")
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "df: %v\n", r)
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

	executeDf(config)
}

