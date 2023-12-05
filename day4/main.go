package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type game struct {
	winning []int
	mine    []int
}

func main() {
	data := loadData()
	ans1 := q1(data)
	ans2 := q2(data)
	fmt.Println("Ans1 is ", ans1)
	fmt.Println("Ans2 is ", ans2)
}

func q1(data []game) int {

	ans := 0
	for _, game := range data {
		match := game.findMatch()

		points := 0
		if match > 0 {
			points = int(math.Pow(2, float64(match-1)))
		}
		ans += points
	}
	return ans
}

func q2(data []game) int {

	ans := 0
	var count []int
	for range data {
		count = append(count, 1)
	}

	for i, game := range data {
		match := game.findMatch()
		for j := 1; j <= match; j++ {
			count[i+j] += count[i]
		}
	}

	for _, v := range count {
		ans += v
	}
	return ans
}

func loadData() []game {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var games []game
	data := strings.Split(string(input), "\n")
	for _, line := range data {
		gameData := strings.Split(line, ":")[1]
		var game game
		numbers := strings.Split(gameData, "|")

		for _, val := range strings.Split(numbers[0], " ") {
			num, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			game.winning = append(game.winning, num)
		}

		for _, val := range strings.Split(numbers[1], " ") {
			num, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			game.mine = append(game.mine, num)
		}
		games = append(games, game)
	}
	return games
}

func (g *game) findMatch() int {
	match := 0
	for _, win := range g.winning {
		for _, mine := range g.mine {
			if win == mine {
				match += 1
				break
			}
		}
	}
	return match
}
