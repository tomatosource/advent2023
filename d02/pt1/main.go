package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	redCap   = 12
	greenCap = 13
	blueCap  = 14
)

type Game struct {
	id     int
	rounds []map[string]int
}

func main() {
	lines := linesFromPath("./adventofcode.com_2023_day_2_input.txt")
	games := gamesFromLines(lines)
	validGames := validGamesFromGames(games)

	var idSum int
	for game := range validGames {
		idSum += game.id
	}
	fmt.Println(idSum)
}

func validGamesFromGames(games <-chan Game) <-chan Game {
	validGames := make(chan Game)

	go func() {
		for game := range games {
			valid := true
			for _, rd := range game.rounds {
				if rd["red"] > redCap ||
					rd["green"] > greenCap ||
					rd["blue"] > blueCap {
					valid = false
					break
				}
			}
			if valid {
				validGames <- game
			}
		}
		close(validGames)
	}()

	return validGames
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
	f, err := os.Open("./adventofcode.com_2023_day_2_input.txt")
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
