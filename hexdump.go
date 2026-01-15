package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	canonical := flag.Bool("C", false, "canonical hex+ASCII")
	length := flag.Uint64("n", 0, "читать N байт")
	skip := flag.Uint64("s", 0, "пропустить N байт")
	help := flag.Bool("h", false, "справка")
	
	flag.Parse()

	if *help { 
		showHelp()
		return 
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "hexdump: missing file operand")
		os.Exit(1)
	}

	for _, filename := range args {
		if err := hexdumpFile(filename, *canonical, *length, *skip); err != nil {
			fmt.Fprintf(os.Stderr, "hexdump: %s: %v\n", filename, err)
		}
		fmt.Println()
	}
}

func hexdumpFile(filename string, canonical bool, maxLen, skip uint64) error {
	f, err := os.Open(filename)
	if err != nil { 
		return err 
	}
	defer f.Close()

	if skip > 0 { 
		if _, err := f.Seek(int64(skip), io.SeekStart); err != nil {
			return fmt.Errorf("seek failed: %v", err)
		}
	}

	totalBytes := uint64(0)
	buf := make([]byte, 16)
	
	for {
		n, err := f.Read(buf)
		if n == 0 { 
			break 
		}
		if err != nil && err != io.EOF { 
			return err 
		}
		
		lineOffset := totalBytes / 16
		printLine(buf[:n], canonical, lineOffset)
		totalBytes += uint64(n)
		
		if maxLen > 0 && totalBytes >= maxLen { 
			break 
		}
	}

	return nil
}

func printLine(data []byte, canonical bool, lineOffset uint64) {
	if canonical {
		fmt.Printf("%08x  ", lineOffset*16)
		
		for i := 0; i < len(data); i++ {
			fmt.Printf("%02x ", data[i])
		}
		for i := len(data); i < 16; i++ {
			fmt.Print("   ")
		}
		
		fmt.Print("  |")
		for _, b := range data {
			if b >= 32 && b <= 126 {
				fmt.Printf("%c", b)
			} else {
				fmt.Print(".")
			}
		}
		for i := len(data); i < 16; i++ {
			fmt.Print(" ")
		}
		fmt.Println("|")
	} else {
		fmt.Printf("%07x ", lineOffset*16)
		for i := 0; i < len(data); i += 2 {
			if i+1 < len(data) {
				val := uint16(data[i]) | uint16(data[i+1])<<8
				fmt.Printf("%04x ", val)
			} else {
				fmt.Printf("%02x ", data[i])
			}
		}
		fmt.Println()
	}
}

func showHelp() {
	fmt.Println("hexdump - отображает содержимое файла в шестнадцатеричном виде")
	fmt.Println()
	fmt.Println("Использование: hexdump [ОПЦИЯ]... ФАЙЛ...")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -C     канонический формат")
	fmt.Println("  -n N   читать N байт")
	fmt.Println("  -s N   пропустить N байт")
	fmt.Println("  -h     показать эту справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  hexdump test.bin")
	fmt.Println("  hexdump -C test.bin              # Канонический формат")
	fmt.Println("  hexdump -s 64 test.bin           # Пропустить первые 64 байта")
	fmt.Println("  hexdump -n 32 test.bin           # Показать первые 32 байта")
}

