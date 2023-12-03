package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const NL = -1

func main() {
	if err := do(); err != nil {
		panic(err)
	}
}

func do() error {
	sum := 0
	leftSet := false
	var left, right int

	for d := range digits() {
		if d == NL {
			val := left*10 + right
			sum += val
			leftSet = false
			left = 0
			right = 0
		} else {
			if leftSet {
				right = d
			} else {
				left = d
				right = d
				leftSet = true
			}
		}
	}

	fmt.Println(sum)

	return nil
}

func digits() <-chan int {
	f, err := os.Open("../pt1/adventofcode.com_2023_day_1_input.txt")
	if err != nil {
		panic(err)
	}

	br := bufio.NewReader(f)
	digits := make(chan int)

	go func() {
		prefixes := []string{}

		for {
			b, err := br.ReadByte()
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				panic(err)
			}

			switch true {
			case b >= '0' && b <= '9':
				digits <- int(b - '0')
				prefixes = []string{}
				break

			case b == '\n':
				digits <- NL
				prefixes = []string{}
				break

			default:
				v := strings.ToLower(string(b))
				prefixes = append(prefixes, "")
				newPrefixes := []string{}

				for _, prefix := range prefixes {
					newPrefix := prefix + v
					digit, validDigit, validPrefix := strToDigit(newPrefix)
					if validDigit {
						digits <- digit
					} else if validPrefix {
						newPrefixes = append(newPrefixes, newPrefix)
					}
				}
				prefixes = newPrefixes
			}
		}

		close(digits)
		f.Close()
	}()

	return digits
}

func strToDigit(s string) (int, bool, bool) {
	validPrefix := false
	for digitVal, digitName := range []string{
		"zero", "one", "two", "three", "four",
		"five", "six", "seven", "eight", "nine",
	} {
		if s == digitName {
			return digitVal, true, false
		} else if strings.HasPrefix(digitName, s) {
			validPrefix = true
		}
	}

	return 0, false, validPrefix
}
