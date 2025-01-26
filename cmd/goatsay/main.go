package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/marufMunshi/goatsay-cli/internal"
)

func main() {
	standardInputFileInfo, standardInputFileError := os.Stdin.Stat()

	if standardInputFileError != nil {
		fmt.Println("Oops!, the Goat encountered on an error", standardInputFileError)
	}

	if standardInputFileInfo.Mode() & os.ModeCharDevice != 0 {
		fmt.Println("Goatsay only work with pipes");
		fmt.Println("Usage: fortune | go run cmd/goatsay/main.go")

		return
	}

	reader := bufio.NewReader(os.Stdin)

	var lines []string

	for {
		inputLine, _, lineReadingError := reader.ReadLine()

		if lineReadingError != nil && lineReadingError == io.EOF {
			break
		}

		lines = append(lines, string(inputLine))
	}

	lines = internal.ConvertTabsToSpaces(lines)
	maxLength, lineIndexToLenghtMap  := internal.CalculateLengthOfLines(lines)
	normalizedLines := internal.NormalizeLinesLength(lines, maxLength, lineIndexToLenghtMap)
	messageToPrint := internal.FormatLinesToBalloonText(normalizedLines, maxLength)

	fmt.Println(messageToPrint)

	// as this file is very small in size, we are reading this file in one step
	// for larger files bufio should be used
	goatFromFile, _ := os.ReadFile("cmd/goatsay/goat.txt")
	fmt.Println(string(goatFromFile))
}