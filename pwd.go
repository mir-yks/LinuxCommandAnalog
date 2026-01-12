package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
)

// Структура для хранения аргументов команды
type PwdOptions struct {
    Help        bool   
    Logical     bool   
    Physical    bool   
    useLogical  bool   
}

// Функция для получения текущей директории с учетом флагов
func getCurrentDir(options PwdOptions) (string, error) {
    if options.useLogical {
        if pwd := os.Getenv("PWD"); pwd != "" {
            if _, err := os.Stat(pwd); err == nil && filepath.IsAbs(pwd) {
                return pwd, nil
            }
        } 
    }
    
    return os.Getwd()
}

// Функция для проверки и обработки флагов
func parseFlags() PwdOptions {
    var options PwdOptions
    
    flag.BoolVar(&options.Help, "h", false, "показать справку об использовании")
    flag.BoolVar(&options.Help, "help", false, "показать справку об использовании")
    flag.BoolVar(&options.Logical, "L", false, "использовать логический путь (разрешить символьные ссылки)")
    flag.BoolVar(&options.Physical, "P", false, "использовать физический путь (без символьных ссылок)")
    
    flag.Parse()
    
    if flag.NArg() > 0 {
        fmt.Fprintf(os.Stderr, "pwd: лишний операнд '%s'\n", flag.Arg(0))
        fmt.Fprintf(os.Stderr, "Используйте 'pwd -h' для получения справки\n")
        os.Exit(1)
    }
    
    options.useLogical = options.Logical && !options.Physical
    
    return options
}

// Функция для отображения справки
func showHelp() {
    fmt.Println("Использование: pwd [КЛЮЧ]...")
    fmt.Println("Выводит полное имя текущей рабочей директории.")
    fmt.Println()
    fmt.Println("Ключи:")
    fmt.Println("  -L             использовать значение PWD из окружения, даже если он содержит символьные ссылки")
    fmt.Println("  -P             использовать физическую структуру каталогов (без символьных ссылок)")
    fmt.Println("  -h, --help     показать эту справку и выйти")
    fmt.Println()
    fmt.Println("Если не указано ни одного ключа, используется -P.")
    fmt.Println()
    fmt.Println("Примеры:")
    fmt.Println("  pwd           # Вывести текущую директорию (эквивалентно pwd -P)")
    fmt.Println("  pwd -L        # Использовать логический путь")
    fmt.Println("  pwd -P        # Использовать физический путь")
}

func main() {
    options := parseFlags()
    
    if options.Help {
        showHelp()
        os.Exit(0)
    }
    
    currentDir, err := getCurrentDir(options)
    if err != nil {
        fmt.Fprintf(os.Stderr, "pwd: ошибка получения текущей директории: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println(currentDir)
}
