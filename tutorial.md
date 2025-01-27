## Checking the input source

First we will look where the input is coming from. As this cli willl only work pipes, we have to check if the input is coming from a pipe or not.

```go
standardInputFileInfo, standardInputFileError := os.Stdin.Stat()
```

Here `Stdin` is the standard input file descriptors. [Stdin doc](https://pkg.go.dev/os#Stdin)
`Stat` funtion returns the file info and error.

If the input is not coming from pipe, we have to show some message and return tahe programe.

```go
	if standardInputFileInfo.Mode() & os.ModeCharDevice != 0 {
		fmt.Println("Goatsay only wants to work with pipe");
		fmt.Println("Usage: fortune | go run cmd/goatsay/main.go ")

		return
	}
```

Here `Mode()` funtion returns binary bits of the input source and `os.ModeCharDevice` has the binary bits of a character device like terminal.

Terminal

```yaml
info.Mode()           = 0000 0000 0000 0000 0000 0000 0001 0000
os.ModeCharDevice     = 0000 0000 0000 0000 0000 0000 0001 0000
-------------------------------------------
Result of `&`         = 0000 0000 0000 0000 0000 0000 0001 0000
```

Pipe

```yaml
info.Mode()           = 0000 0000 0000 0000 0000 0000 0000 0001
os.ModeCharDevice     = 0000 0000 0000 0000 0000 0000 0001 0000
-------------------------------------------
Result of `&`         = 0000 0000 0000 0000 0000 0000 0000 0000
```

Bitwise operator reference -> https://stackoverflow.com/questions/3427585/understanding-the-bitwise-and-operator

## Read lines form input and build an array of strings

Now we will read the input source `os.Stdin`. We will use [bufio](https://pkg.go.dev/bufio) to read the input in a buffered form.

It reads data in memory in chunks (a buffer) instead of byte-by-byte or line-by-line directly from the input source. So, it minimizes the number of I/O operations by reading a larger chunk of data (a buffer) into memory at once. [See more](https://www.educative.io/answers/how-to-read-and-write-with-golang-bufio)

```go
	reader := bufio.NewReader(os.Stdin)
```

Here `reader` variable is a pointer to the bufio reader struct and has access to all the methods. We will only need `ReadLine()` method.
Go is not a OOP language, it does not have classes. So, worth to understand that how this grouping is done. In go, we can write fuctions with a special reciver type and those will become method to the type and can be accessed with dot notation. [Read this blog to know more](https://golangbot.com/methods/)

So, after reading input from line by line we will appened those in a slice of strings.

```go
	reader := bufio.NewReader(os.Stdin)

	var lines []string

	for {
		inputLine, _, lineReadingError := reader.ReadLine()

		if lineReadingError != nil && lineReadingError == io.EOF {
			break
		}

		lines = append(lines, string(inputLine))
	}
```

## Convert tabs to space for counting runes correctly

We already have the input lines in a string slice, now we need to replace any tab ("/t") to 4 spaces. It will help us count runes correctly.

[Replace doc](https://pkg.go.dev/strings#Replace), [append doc](https://pkg.go.dev/builtin#append)

```go
func ConvertTabsToSpaces(lines []string) []string {
	var result []string

	for _, line := range lines {
		line = strings.Replace(line, "\t", "    ", -1)
		result = append(result, line)
	}

	return result
}
```

## Calculate maximum length of lines

We need to calculate the length of the lines and the maximum length. This will help us formatted all the text nicely.
We will be using `utf8.RuneCountInString` function insteasd of `len`. `len` returns byte count, so if a string contains special character it will not return the correct character count.
[RuneCountInString example](https://pkg.go.dev/unicode/utf8#RuneCountInString)

One thing worth mentiong is map datatype is passed as reference by default. We only need one copy of the `lineIndexToLenghtMap` through out the application

```go
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
```

## Normalizing all the lines

We already know the maximum length and all the length of lines. So, this function will add extra space to make all the length equal. It will help to wrap the text inside a ballon.

```go
func NormalizeLinesLength(lines []string, maximumLength int, lineIndexToLenghtMap map[int]int ) []string {
	var result []string;

	for index, line := range lines {
		line = line + strings.Repeat(" ", maximumLength - lineIndexToLenghtMap[index])

		result = append(result, line)
	}

	return result
}
```

## Format lines to balllon text

In this section we will just format the lines of string to a single string in a ballon

```go
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
```

## Print the Goat

Last part is to printing the goat. I have taken this goat ASCII art form [arcii art archive website](https://www.asciiart.eu/animals/other-land). This Goat art contains bacticks (``) and other special characters. In go we can not use backtick inside backtick. so, it best to be used form a text file.

```go
	goatFromFile, _ := os.ReadFile("cmd/goatsay/goat.txt")
```

Initially I read the file using a realtive path. It works while running with `go run`. But when this program installed with `go install` and run as a binary
working directory is different and program fail to find the .txt file.

So, mitigate this problem we used embed directive to load the file at complie time into a variable. [embed directive](https://pkg.go.dev/embed#hdr-Directives)

```go
//go:embed goat.txt
var goatArt string
```
