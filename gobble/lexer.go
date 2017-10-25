// Lexer divides a given Scribble file into a slice of string tokens

// The lexer reads in a local Scribble protocol from file and processes
// it into an array of lexical tokens. These tokens consist of strings
// which can be either runs of alphanumeric characters such as "protocol"
// or punctuation characters with specific meanings in Scribble, such as
// "{" or ";".

// input: validated Scribble .scr local protocol file

// output: []string

package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func ReadLines(inputFile string) []string {
	lines := make([]string, 0)
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "	 ")
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func RemoveLineComments(lines []string) string {
	outputString := ""
	for _, line := range lines {
		for pos, char := range line {
			if pos < len(line)-1 {
				if char == '/' && line[pos+1] == '/' {
					break
				} else {
					outputString += string(char)
				}
			} else {
				outputString += string(char)
			}
		}
		outputString += string(' ')
	}
	return outputString
}

func RemoveBlockComments(in string) string {
	out := ""
	inBlockComment := false
	var prevChar rune
	for _, char := range in {
		if prevChar == '/' && char == '*' {
			inBlockComment = true
		}
		if !inBlockComment && char != '/' {
			out += string(char)
		}
		if prevChar == '*' && char == '/' {
			inBlockComment = false
		}
		prevChar = char
	}
	return out
}

func isSpecialChar(r rune) bool {
	specialChars := []rune{'<', '>', '/', '"', '.', ',', '\\', ';', '{', '}', '[', ']', '(', ')'}
	isSpecial := false
	for _, currentRune := range specialChars {
		if currentRune == r {
			isSpecial = true
		}
	}
	return isSpecial
}

func Tokenize(words []string) []string {
	tokens := make([]string, 0)
	currentWord := ""
	for _, word := range words {
		for _, char := range word {
			if isSpecialChar(char) {
				tokens = append(tokens, currentWord)
				tokens = append(tokens, string(char))
				currentWord = ""
			} else {
				currentWord += string(char)
			}
		}
		tokens = append(tokens, currentWord)
		currentWord = ""
	}
	return tokens
}

func RemoveEmptyStrings(tokensWithEmptyStrings []string) []string {
	tokensWithoutEmptyStrings := make([]string, 0)
	for _, token := range tokensWithEmptyStrings {
		if token != "" {
			tokensWithoutEmptyStrings = append(tokensWithoutEmptyStrings, token)
		}
	}
	return tokensWithoutEmptyStrings
}

func GetTokens(inputFile string) []string {
	lines := ReadLines(inputFile)
	linesAsSingleString := RemoveLineComments(lines)
	linesAsSingleString = RemoveBlockComments(linesAsSingleString)
	tokensDividedByWhiteSpace := strings.Fields(linesAsSingleString)
	tokens := Tokenize(tokensDividedByWhiteSpace)
	tokens = RemoveEmptyStrings(tokens)
	return tokens
}
