package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	one()
	two()
}

func one() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	answer, gameNo := 0, 0
	for scanner.Scan() {
		gameNo++
		text := scanner.Text()
		list := strings.Split(text, ": ")[1]
		draws := strings.Split(list, "; ")
		failed := false
		for _, draw := range draws {
			ball := strings.Split(draw, ", ")
			for i := 0; i < len(ball); i++ {
				ballparts := strings.Split(ball[i], " ")
				val, _ := strconv.ParseInt(ballparts[0], 10, 8)
				color := ballparts[1]
				// fmt.Println("color and val", color, val)
				if color == "red" && val > 12 {
					failed = true
					break
				} else if color == "green" && val > 13 {
					failed = true
					break
				} else if color == "blue" && val > 14 {
					failed = true
					break
				}

			}
		}
		if !failed {
			answer += gameNo
		} else {
		}
	}
	fmt.Println("Part 1 Answer: ", answer)
}

func two() {

	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	answer := 0
	for scanner.Scan() {
		text := scanner.Text()
		list := strings.Split(text, ": ")[1]
		draws := strings.Split(list, "; ")
		maxRed, maxGreen, maxBlue := 0, 0, 0
		for _, draw := range draws {
			ball := strings.Split(draw, ", ")
			for i := 0; i < len(ball); i++ {
				ballparts := strings.Split(ball[i], " ")
				val, _ := strconv.ParseInt(ballparts[0], 10, 8)
				color := ballparts[1]
				if color == "red" && val > int64(maxRed) {
					maxRed = int(val)
				} else if color == "green" && val > int64(maxGreen) {
					maxGreen = int(val)
				} else if color == "blue" && val > int64(maxBlue) {
					maxBlue = int(val)
				}

			}
		}
		linePower := maxBlue * maxGreen * maxRed
		// fmt.Println("linePower: ", linePower)
		answer += linePower
	}
	fmt.Println("Part 2 Answer: ", answer)
}
