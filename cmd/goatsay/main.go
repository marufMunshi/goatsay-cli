package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/marufMunshi/goatsay-cli/internal"
)

//go:embed goat.txt
var goatArt string

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
	fmt.Println(goatArt)
}