package internal

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// this function will convert tab to space by replacing
// "\t" with 4 space
// it will help counting runes correctly to prevent misalignments
func ConvertTabsToSpaces(lines []string) []string {
	var result []string

	for _, line := range lines {
		line = strings.Replace(line, "\t", "    ", -1)
		result = append(result, line)
	}

	return result
}

func CalculateLengthOfLines(lines []string) (int, map[int]int) {
	maxLength := 0
	lineIndexToLenghtMap := make(map[int]int)

	for index, line := range lines {
		currentLineLength := utf8.RuneCountInString(line)

		lineIndexToLenghtMap[index] = currentLineLength;

		if currentLineLength > maxLength {
			maxLength = currentLineLength
		}
	}

	return maxLength, lineIndexToLenghtMap;
}

func NormalizeLinesLength(lines []string, maximumLength int, lineIndexToLenghtMap map[int]int ) []string {
	var result []string;

	for index, line := range lines {
		line = line + strings.Repeat(" ", maximumLength - lineIndexToLenghtMap[index])

		result = append(result, line)
	}

	return result
}

func FormatLinesToBalloonText(lines []string, maximumLength int) string {
	var result []string

	topBorder := " " + strings.Repeat("_", maximumLength + 2)
	bottomBorder := " " + strings.Repeat("-", maximumLength + 2)

	result = append(result, topBorder)

	linesCount := len(lines)

	if linesCount == 1 {
		singleLine := fmt.Sprintf("%s %s %s", "<", lines[0], ">")
		result = append(result, singleLine)
		result = append(result, bottomBorder)

		return strings.Join(result, "\n")
	}

	for index, line := range lines {
		fromattedString := "";

		if index == 0 {
			fromattedString = fmt.Sprintf("%s %s %s", "/", line, "\\")
		} else if index == linesCount - 1 {
			fromattedString = fmt.Sprintf("%s %s %s", "\\", line, "/")
		} else {
			fromattedString = fmt.Sprintf("%s %s %s", "|", line, "|")
		}

		result = append(result, fromattedString)
	}

	result = append(result, bottomBorder)

	return strings.Join(result, "\n")
}