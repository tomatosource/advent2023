package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Game struct {
	id     int
	rounds []map[string]int
}

func main() {
	lines := linesFromPath("../pt1/adventofcode.com_2023_day_2_input.txt")
	games := gamesFromLines(lines)
	powers := powersFromGames(games)

	var sum int
	for power := range powers {
		sum += power
	}
	fmt.Println(sum)
}

func powersFromGames(games <-chan Game) <-chan int {
	powers := make(chan int)

	go func() {
		for game := range games {
			minRed, minGreen, minBlue := 0, 0, 0
			for _, rd := range game.rounds {
				if red := rd["red"]; red > minRed {
					minRed = red
				}
				if blue := rd["blue"]; blue > minBlue {
					minBlue = blue
				}
				if green := rd["green"]; green > minGreen {
					minGreen = green
				}
			}
			powers <- minRed * minBlue * minGreen
		}

		close(powers)
	}()

	return powers
}

func gamesFromLines(lines <-chan string) <-chan Game {
	games := make(chan Game)

	go func() {
		for line := range lines {
			colonIdx := strings.Index(line, ":")
			id, _ := strconv.Atoi(line[5:colonIdx])
			rounds := []map[string]int{}

			for _, roundStr := range strings.Split(line[colonIdx+2:], "; ") {
				r := map[string]int{}
				for _, pickStr := range strings.Split(roundStr, ", ") {
					parts := strings.Split(pickStr, " ")
					count, _ := strconv.Atoi(parts[0])
					r[parts[1]] = count
				}
				rounds = append(rounds, r)
			}

			games <- Game{id, rounds}
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
