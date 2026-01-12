package main

import (
 "flag"
 "fmt"
 "os"
)

func printCurrentDir() (string, error) {
 return os.Getwd()
}

func main() {
 help := flag.Bool("h", false, "показать помощь")
 flag.Parse()

 if *help {
  fmt.Println("Использование: pwd [-h]")
  fmt.Println("-h: помощь")
  os.Exit(0)
 }

 currentDir, err := printCurrentDir()
 if err != nil {
  fmt.Println("Ошибка:", err)
  os.Exit(1)
 }

 fmt.Println(currentDir)
}
