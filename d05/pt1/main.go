package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type adapter struct {
	source        string
	destination   string
	rangeAdapters []rangeAdapter
}

type rangeAdapter struct {
	destinationRangeStart int
	sourceRangeStart      int
	rangeLength           int
}

func main() {
	adapters, seeds := readInput("./adventofcode.com_2023_day_5_input.txt")

	fmt.Println(adapters["humidity"])

	min := math.MaxInt

	for _, seed := range seeds {
		val := seed
		for out := "seed"; out != "location"; {
			val, out = adapters[out].convert(val)
		}
		if val < min {
			min = val
		}
	}

	fmt.Println(min)
}

func readInput(path string) (map[string]adapter, []int) {
	file, _ := os.Open(path)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	seedLine := strings.TrimPrefix(lines[0], "seeds: ")

	var seeds []int
	substrings := strings.Split(seedLine, " ")
	for _, s := range substrings {
		num, _ := strconv.Atoi(s)
		seeds = append(seeds, num)
	}

	adapters := map[string]adapter{}
	a := adapter{}
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		if strings.Contains(line, "map:") {
			adapters[a.source] = a

			parts := strings.Split(strings.TrimSuffix(line, " map:"), "-")
			a = adapter{
				source:        parts[0],
				destination:   parts[2],
				rangeAdapters: []rangeAdapter{},
			}
		} else {
			parts := strings.Split(line, " ")
			destinationRangeStart, _ := strconv.Atoi(parts[0])
			sourceRangeStart, _ := strconv.Atoi(parts[1])
			rangeLength, _ := strconv.Atoi(parts[2])
			a.rangeAdapters = append(a.rangeAdapters, rangeAdapter{
				destinationRangeStart,
				sourceRangeStart,
				rangeLength,
			})
		}
	}
	adapters[a.source] = a

	return adapters, seeds
}

func (a adapter) convert(s int) (int, string) {
	for _, ra := range a.rangeAdapters {
		if s >= ra.sourceRangeStart && s < ra.sourceRangeStart+ra.rangeLength {
			return ra.destinationRangeStart + (s - ra.sourceRangeStart), a.destination
		}
	}

	return s, a.destination
}
