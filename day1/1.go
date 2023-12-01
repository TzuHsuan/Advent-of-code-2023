package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	two()
}

func one() {
	digits := "1234567890"
	file, err := os.Open("input1.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	sum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		firstIndex := strings.IndexAny(text, digits)
		lastIndex := strings.LastIndexAny(text, digits)
		numString := string([]byte{text[firstIndex], text[lastIndex]})
		number, err := strconv.ParseInt(numString, 10, 64)
		if err != nil {
			panic(err)
		}
		fmt.Println(text)
		sum += int(number)
	}

	fmt.Println(sum)
}

func two() {
	file, err := os.Open("input1.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	replace := []struct {
		str    string
		number string
	}{
		{"one", "1"},
		{"two", "2"},
		{"three", "3"},
		{"four", "4"},
		{"five", "5"},
		{"six", "6"},
		{"seven", "7"},
		{"eight", "8"},
		{"nine", "9"},
		{"zero", "0"},
		{"1", "1"},
		{"2", "2"},
		{"3", "3"},
		{"4", "4"},
		{"5", "5"},
		{"6", "6"},
		{"7", "7"},
		{"8", "8"},
		{"9", "9"},
		{"0", "0"},
	}
	sum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		left := findLeft(text, replace)
		right := findRight(text, replace)

		number, err := strconv.ParseInt(left+right, 10, 64)
		if err != nil {
			panic(err)
		}
		sum += int(number)
	}

	fmt.Println(sum)
}

func findLeft(str string, table []struct {
	str    string
	number string
}) string {
	for {
		for _, v := range table {
			if strings.HasPrefix(str, v.str) {
				return v.number
			}
		}
		str = str[1:]
	}
}

func findRight(str string, table []struct {
	str    string
	number string
}) string {
	for {
		for _, v := range table {
			if strings.HasSuffix(str, v.str) {
				return v.number
			}
		}
		str = str[:len(str)-1]
	}
}
