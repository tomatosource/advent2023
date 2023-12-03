package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := do(); err != nil {
		panic(err)
	}
}

func do() error {
	f, err := os.Open("./adventofcode.com_2023_day_1_input.txt")
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	br := bufio.NewReader(f)
	sum := 0
	leftSet := false
	var left, right byte

	for {
		b, err := br.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("reading byte: %w", err)
		}

		switch true {
		case b >= '0' && b <= '9':
			if leftSet {
				right = b
			} else {
				left = b
				right = b
				leftSet = true
			}
			break

		case b == '\n':
			sum += int((left-'0')*10 + (right - '0'))
			leftSet = false
			break
		}
	}

	fmt.Println(sum)

	return nil
}
