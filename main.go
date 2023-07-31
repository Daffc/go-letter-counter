package main

import (
	"fmt"
	"os"

	utils "github.com/Daffc/go-letter-counter/package/utils"
)

func main() {
	nGoroutines, bufferSize, fInput, fOutput := utils.HandleArguemnts(os.Args[1:])
	defer fInput.Close()
	defer fOutput.Close()

	fmt.Println(nGoroutines, bufferSize, fInput, fOutput)
}
