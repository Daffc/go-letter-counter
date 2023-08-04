package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	utils "github.com/Daffc/go-letter-counter/package/utils"
)

func BufferedWorker(id int, wg *sync.WaitGroup, fInput *os.File, bufferSize int, result []int) {
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
			if 65 <= buffer[i] && buffer[i] <= 90 {
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

func NotBufferedWorker(id int, wg *sync.WaitGroup, inputBuffer []byte, startPosition int, maxPosition int, result []int) {
	for i := startPosition; i < maxPosition; i++ {

		if 65 <= inputBuffer[i+startPosition] && inputBuffer[i+startPosition] <= 90 {
			result[inputBuffer[i+startPosition]-65]++
		} else {
			if 97 <= inputBuffer[i+startPosition] && inputBuffer[i+startPosition] <= 122 {
				result[inputBuffer[i+startPosition]-97]++
			}
		}

	}

	wg.Done()
}

func main() {
	nGoroutines, bufferSize, fInput, fOutput, isBuffered := utils.HandleArguemnts(os.Args[1:])
	defer fInput.Close()
	defer fOutput.Close()

	var waitingGroup sync.WaitGroup

	resultBuffers := make([][]int, nGoroutines)
	finalResult := make([]int, 26)

	for i := 0; i < nGoroutines; i++ {
		resultBuffers[i] = make([]int, 26)
	}

	if isBuffered {
		for i := 0; i < nGoroutines; i++ {
			waitingGroup.Add(1)
			go BufferedWorker(i, &waitingGroup, fInput, bufferSize, resultBuffers[i])
		}
	} else {

		// Readding and buffering whole file.
		bufferReadFile, err := os.ReadFile(fInput.Name())
		utils.CheckError(err)

		sliceSize := len(bufferReadFile) / nGoroutines
		reminder := len(bufferReadFile) % nGoroutines

		// Initiating all go routines except the last one.
		for i := 0; i < nGoroutines-1; i++ {
			waitingGroup.Add(1)
			go NotBufferedWorker(i, &waitingGroup, bufferReadFile, i*sliceSize, sliceSize, resultBuffers[i])
		}

		// Initiating the last goroutine adding remind to the end of the read buffer.
		waitingGroup.Add(1)
		go NotBufferedWorker(nGoroutines-1, &waitingGroup, bufferReadFile, (nGoroutines-1)*sliceSize, sliceSize+reminder, resultBuffers[nGoroutines-1])
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
