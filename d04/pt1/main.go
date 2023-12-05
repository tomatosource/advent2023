package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Game struct {
	id             string
	winningNumbers []int
	numbers        []int
	line           string
}

var cardPattern = regexp.MustCompile(`Card\s+(\d+): ([\d\s]+) \| ([\d\s]+)`)

func main() {
	lines := linesFromPath("./adventofcode.com_2023_day_4_input.txt")
	games := gamesFromLines(lines)

	var sum int
	for g := range games {
		fmt.Printf("%d - %s\n\n", g.Score(), g.line)
		sum += g.Score()
	}
	fmt.Println(sum)
}

func gamesFromLines(lines <-chan string) <-chan Game {
	games := make(chan Game)

	go func() {
		for line := range lines {
			fmt.Println(line)
			matches := cardPattern.FindStringSubmatch(line)
			id := matches[1]
			winningNumbers := parseIntList(matches[2])
			numbers := parseIntList(matches[3])

			fmt.Println(line)
			games <- Game{
				id, winningNumbers, numbers, line,
			}
		}
		close(games)
	}()

	return games
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

func (g *Game) Score() int {
	matches := 0
	for _, x := range g.winningNumbers {
		for _, y := range g.numbers {
			if x == y {
				matches++
				break
			}
		}
	}

	if matches == 0 {
		return 0
	}

	return 1 << (matches - 1)
}

func parseIntList(input string) []int {
	parts := strings.Fields(input)
	var result []int
	for _, part := range parts {
		num, _ := strconv.Atoi(part)
		result = append(result, num)
	}
	fmt.Println(result)
	return result
}
