package lettercounter

import (
	"flag"
	"os"
)

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

/*
Handles arguments, returning:

	number of Worker goroutines (int).
	buffer size for each reader (int).
	input file descriptor (*os.File)
	output file descriptor (*os.File)
*/
func HandleArguemnts(args []string) (int, int, *os.File, *os.File, bool) {
	inputFilePathPtr := flag.String("i", "", "PathPtr to the file to be read.")
	outputFilePathPtr := flag.String("o", "", "PathPtr to the output file.")
	nGoroutines := flag.Int("n", 1, "Number of worker goroutines.")
	bufferSize := flag.Int("b", 2000, "Buffer size (bytes) for each worker.")
	isBuffered := flag.Bool("buffered", false, "Flag that indicates if the reads will be buffered for each worker (if not set there will be only one buffer for the whole read shared between the workers).")
	flag.Parse()

	var inputFile *os.File
	var outputFile *os.File
	var err error

	if *inputFilePathPtr == "" {
		inputFile = os.Stdin
	} else {
		inputFile, err = os.OpenFile(*inputFilePathPtr, os.O_RDONLY, 0444)
		CheckError(err)
	}

	if *outputFilePathPtr == "" {
		outputFile = os.Stdout
	} else {
		outputFile, err = os.OpenFile(*outputFilePathPtr, os.O_CREATE|os.O_WRONLY, 0666)
		CheckError(err)
	}

	return *nGoroutines, *bufferSize, inputFile, outputFile, *isBuffered
}
