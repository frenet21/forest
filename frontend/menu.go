package main

import (
	"fmt"
	"bufio"
	"os"
)
func main(){
	fmt.Printf("★★★ Welcome to Forest! You can begin your chatting now ★★★")
	fmt.Printf("//(Enter any number of options above)")

	//initial a bufio
	reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter text: ")
    text, _ := reader.ReadString('\n')
    fmt.Println(text)
}