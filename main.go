package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	utils "github.com/Daffc/go-letter-counter/package/utils"
)

func worker(id int, wg *sync.WaitGroup, fInput *os.File, bufferSize int, result []int) {
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
			if 65 <= buffer[i] && buffer[i] <= 95 {
				result[buffer[i]-65]++
			} else {
				if 97 <= buffer[i] && buffer[i] <= 122 {
					result[buffer[i]-97]++
				}
			}

		}
	}

	wg.Done()
}

func main() {
	nGoroutines, bufferSize, fInput, fOutput := utils.HandleArguemnts(os.Args[1:])
	defer fInput.Close()
	defer fOutput.Close()

	var waitingGroup sync.WaitGroup

	resultBuffers := make([][]int, nGoroutines)
	finalResult := make([]int, 26)

	for i := 0; i < nGoroutines; i++ {
		resultBuffers[i] = make([]int, 26)
	}

	for i := 0; i < nGoroutines; i++ {
		waitingGroup.Add(1)
		go worker(i, &waitingGroup, fInput, bufferSize, resultBuffers[i])
	}

	waitingGroup.Wait()

	for i := 0; i < nGoroutines; i++ {
		for j := 0; j < 26; j++ {
			finalResult[j] += resultBuffers[i][j]
		}
	}

	for i := 0; i < 26; i++ {
		fmt.Fprintf(fOutput, "%c = %d\n", i+65, finalResult[i])
	}

}
