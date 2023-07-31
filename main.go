package main

import (
	"fmt"
	"io"
	"os"

	utils "github.com/Daffc/go-letter-counter/package/utils"
)

func main() {
	nGoroutines, bufferSize, fInput, fOutput := utils.HandleArguemnts(os.Args[1:])
	defer fInput.Close()
	defer fOutput.Close()

	buffer := make([]byte, bufferSize)

	for {
		readBytes, err := fInput.Read(buffer)

		if err != nil {
			if err != io.EOF {
				panic(err)
			}

			break
		}

		for i := 0; i < readBytes; i++ {
			fmt.Fprint(fOutput, string(buffer[i]))
		}
	}

	fmt.Println(nGoroutines, bufferSize, fInput, fOutput)
}
