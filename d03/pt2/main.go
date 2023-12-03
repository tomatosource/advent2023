package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

const Symbol = -1

type Token struct {
	start  int
	end    int
	row    int
	number int
	ratio  int
}

func (t *Token) adjacentPart(as Token) bool {
	if as.number == Symbol {
		return false
	}

	if t.row == as.row {
		return as.end+1 == t.start || t.end+1 == as.start
	} else {
		return (as.start >= t.start-1 && as.start <= t.end+1) ||
			(as.end >= t.start-1 && as.start <= t.end+1)
	}
}

func main() {
	lines := linesFromPath("../pt1/adventofcode.com_2023_day_3_input.txt")
	tokens := tokensFromLines(lines)
	grid := gridFromTokens(tokens)
	gears := gearsFromGrid(grid)

	var sum int
	for gear := range gears {
		sum += gear.ratio
	}
	fmt.Println(sum)
}

func gearsFromGrid(grid [][]Token) <-chan Token {
	gears := make(chan Token)

	go func() {
		for rowIdx, row := range grid {
			for tokenIdx, token := range row {
				if token.number != Symbol {
					continue
				}

				aps := []Token{}

				if tokenIdx > 0 && token.adjacentPart(row[tokenIdx-1]) {
					aps = append(aps, row[tokenIdx-1])
				}

				// right
				if tokenIdx < len(row)-1 && token.adjacentPart(row[tokenIdx+1]) {
					aps = append(aps, row[tokenIdx+1])
				}

				// row above
				if rowIdx > 0 {
					for _, ap := range grid[rowIdx-1] {
						if token.adjacentPart(ap) {
							aps = append(aps, ap)
						} else if ap.start > token.end {
							break
						}
					}
				}

				// row below
				if rowIdx < len(grid)-1 {
					for _, ap := range grid[rowIdx+1] {
						if token.adjacentPart(ap) {
							aps = append(aps, ap)
						} else if ap.start > token.end {
							break
						}
					}
				}

				if len(aps) == 2 {
					token.ratio = aps[0].number * aps[1].number
					gears <- token
				}
			}
		}
		close(gears)
	}()

	return gears
}

func gridFromTokens(tokens <-chan Token) [][]Token {
	grid := [][]Token{}
	rc := 0
	r := []Token{}

	for token := range tokens {
		if rc != token.row {
			rc = token.row
			grid = append(grid, r)
			r = []Token{}
		}
		r = append(r, token)
	}
	grid = append(grid, r)

	return grid
}

func tokensFromLines(lines <-chan string) <-chan Token {
	tokens := make(chan Token)

	go func() {
		row := 0
		for line := range lines {
			lineBytes := []byte(line)
			inNumber := false
			var start int

			for i, c := range lineBytes {
				switch true {
				case c == '.':
					continue

				case !isDigit(c):
					tokens <- Token{
						start:  i,
						end:    i,
						row:    row,
						number: Symbol,
					}
					break

				default:
					if !inNumber {
						start = i
						inNumber = true
					}

					if i >= len(lineBytes)-1 || !isDigit(lineBytes[i+1]) {
						number, _ := strconv.Atoi(line[start : i+1])
						end := i
						tokens <- Token{
							start:  start,
							end:    end,
							row:    row,
							number: number,
						}
						inNumber = false
					}
				}
			}
			row++
		}
		close(tokens)
	}()

	return tokens
}

func linesFromPath(path string) <-chan string {
	lines := make(chan string)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)

	go func() {
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			lines <- string(line)
		}
		f.Close()
		close(lines)
	}()

	return lines
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (t Token) String() string {
	return fmt.Sprintf(
		"start/end: %3d/%3d - row %3d - number %3d",
		t.start, t.end,
		t.row, t.number,
	)
}
