package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	// one(file)
	two(file)
	defer file.Close()
}

func one(file *os.File) {
	scanner := bufio.NewScanner(file)
	data := new([140]string)

	// wait := bufio.NewScanner(os.Stdin)

	line := 0
	for scanner.Scan() {
		data[line] = scanner.Text()
		line++
	}

	for _, dat := range data {
		fmt.Println(dat)
	}

	answer := 0
	for lineNo, str := range data {
		numStr := ""
		for i := 0; i < len(str); i++ {
			numStr += string(data[lineNo][i])
			if _, err := strconv.Atoi(numStr); err != nil {
				if len(numStr) > 1 { //End of number sequence
					numStr = numStr[:len(numStr)-1]
					if findSymbol(i, lineNo, len(numStr), data) {
						num, _ := strconv.Atoi(numStr)
						answer += num
						replacement := ""
						for range numStr {
							replacement += "."
						}

						data[lineNo] = data[lineNo][:i-len(numStr)] + replacement + data[lineNo][i:]
					}
				}
				numStr = ""
				continue
			}
		}

		if len(numStr) > 1 { //End of number sequence
			if findSymbol(140, lineNo, len(numStr), data) {
				num, _ := strconv.Atoi(numStr)
				fmt.Println(num)
				answer += num
				replacement := ""
				for range numStr {
					replacement += "."
				}

				data[lineNo] = data[lineNo][:140-len(numStr)] + replacement + data[lineNo][140:]
			}
		}
	}
	for _, dat := range data {
		fmt.Println(dat)
	}
	fmt.Println("answer is ", answer)
}

func findSymbol(x int, y int, length int, data *[140]string) bool {
	symbolChars := "!@#$%^&*()_+-=/"

	testStr := ""
	for row := y - 1; row < y+2; row++ {
		if row < 0 || row > 139 {
			continue
		}
		for col := x - length - 1; col < x+1; col++ {
			if col < 0 || col > 139 {
				continue
			}
			testStr += string((*data)[row][col])
		}
	}
	// fmt.Println(testStr)
	// if strings.Contains(testStr, "193") {
	// 	fmt.Println(testStr)
	// }
	return strings.ContainsAny(testStr, symbolChars)
}

func two(file *os.File) {
	scanner := bufio.NewScanner(file)
	line := 0
	numbers := new([140][140]int)
	gearChar := "*"
	numberChars := "1234567890"
	type position struct {
		x int
		y int
	}
	symbols := make([]position, 0, 200)
	answer := 0
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		prevNumChar := ""
		for i := 0; i < len(scanner.Text()); i++ {
			if strings.ContainsAny(string(scanner.Text()[i]), numberChars) {
				prevNumChar += string(scanner.Text()[i])
			} else {
				if len(prevNumChar) != 0 {
					parsed, _ := strconv.ParseInt(prevNumChar, 10, 0)
					for h := 0; h < len(prevNumChar); h++ {
						numbers[line][i-h-1] = int(parsed)
					}
				}
				prevNumChar = ""
			}

			if scanner.Text()[i] == '.' {
				continue
			} else if strings.Contains(string(scanner.Text()[i]), gearChar) {
				symbols = append(symbols, position{x: i, y: line})
			}
		}
		if len(prevNumChar) != 0 {
			parsed, _ := strconv.ParseInt(prevNumChar, 10, 0)
			for h := 0; h < len(prevNumChar); h++ {
				numbers[line][140-h-1] = int(parsed)
			}
		}
		line++
	}
	// fmt.Println(numbers)
	fmt.Println(symbols)

	for _, v := range symbols {
		prevNum := false
		adj := make([]int, 0, 9)
		if v.y != 0 {
			row := v.y - 1
			if v.x != 0 {
				if numbers[row][v.x-1] != 0 {
					// answer += numbers[row][v.x-1]
					adj = append(adj, numbers[row][v.x-1])
					prevNum = true
				}
			}

			if numbers[row][v.x] != 0 && prevNum == false {
				// answer += numbers[row][v.x]
				adj = append(adj, numbers[row][v.x])
				prevNum = true
			}

			if numbers[row][v.x] == 0 && prevNum == true {
				prevNum = false
			}

			if prevNum == false && v.x != 140 && numbers[row][v.x+1] != 0 {
				// answer += numbers[row][v.x+1]
				adj = append(adj, numbers[row][v.x+1])
			}
		}
		row := v.y
		prevNum = false

		if v.x != 0 && numbers[row][v.x-1] != 0 {
			// answer += numbers[row][v.x-1]
			adj = append(adj, numbers[row][v.x-1])
		}
		if v.x != 140 && numbers[row][v.x+1] != 0 {
			// answer += numbers[row][v.x+1]
			adj = append(adj, numbers[row][v.x+1])
		}
		if v.y != 140 {
			row := v.y + 1
			if v.x != 0 {
				if numbers[row][v.x-1] != 0 {
					// answer += numbers[row][v.x-1]
					adj = append(adj, numbers[row][v.x-1])
					prevNum = true
				}
			}

			if numbers[row][v.x] != 0 && prevNum == false {
				// answer += numbers[row][v.x]
				adj = append(adj, numbers[row][v.x])
				prevNum = true
			}

			if numbers[row][v.x] == 0 && prevNum == true {
				prevNum = false
			}

			if prevNum == false && v.x != 140 && numbers[row][v.x+1] != 0 {
				adj = append(adj, numbers[row][v.x+1])
				// answer += numbers[row][v.x+1]
			}
		}
		fmt.Println(adj)
		if len(adj) > 1 {
			ratio := 1
			for _, val := range adj {
				ratio *= val
			}
			answer += ratio
		}
	}
	fmt.Println("Answer 2 is ", answer)
}
