package main

import (
	"cmp"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type pod struct {
	stage string
	value int
}

type podRange struct {
	stage    string
	valStart int
	valEnd   int
}

type transformer struct {
	inputType  string
	outputType string
	chart      []Chart
}

type Chart struct {
	begin int
	end   int
	diff  int
}

func main() {
	transformers, pods := loadData(1)
	var answer1 int
	answer1 = findLowestLocation(pods, transformers)
	fmt.Println("Answer 1", answer1)
	t2, p2 := loadData(2)
	var answer2 int
	answer2 = findLowestLocation(p2, t2)
	fmt.Println("Answer 2", answer2)
}

func findLowestLocation(pods []podRange, trans map[string]transformer) int {
	answer := 0
	for len(pods) != 0 {
		for pods[0].stage != "location" {
			currentCharts := trans[pods[0].stage].chart
			for cp, pp := 0, pods[0].valStart; pp <= pods[0].valEnd; {
				cc := currentCharts[cp]
				if pp > cc.end {
					cp += 1
					if cp == len(currentCharts) {
						pods[0].stage = trans[pods[0].stage].outputType
						break
					}
					continue
				}

				np := podRange{}
				//Starts outside of map, No change is value
				if pp < cc.begin {
					np.valStart = pp
					//Entire range is outside map
					if pods[0].valEnd < cc.begin {
						np.valEnd = pods[0].valEnd
						np.stage = trans[pods[0].stage].outputType
						pods[0] = np
						break
					} else { //Extends into map, create split
						np = podRange{valStart: pp, valEnd: cc.begin - 1, stage: trans[pods[0].stage].outputType}
						pods = append(pods, np)
						pods[0].valStart = cc.begin
						pp = cc.begin
					}
				} else { //Starts in map
					//Entire range is in map
					if pods[0].valEnd <= cc.end {
						np = podRange{valStart: pp - cc.diff, valEnd: pods[0].valEnd - cc.diff, stage: trans[pods[0].stage].outputType}
						pods[0] = np
						break
					} else { //Extends out of map, create split
						np = podRange{valStart: pp - cc.diff, valEnd: cc.end - cc.diff, stage: trans[pods[0].stage].outputType}
						pods = append(pods, np)
						pods[0].valStart = cc.end + 1
						pp = cc.end + 1
					}
				}
			}
		}
		if pods[0].valStart < answer || answer == 0 {
			answer = pods[0].valStart
		}
		pods = pods[1:]
	}
	return answer
}

func loadData(question int) (map[string]transformer, []podRange) {
	file, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")

	var pods []podRange
	start := strings.Split(lines[0], " ")
	startType := start[0][:len(start[0])-2]
	start = start[1:]
	if question == 1 {

		for i := 0; i < len(start); i++ {
			valStart, _ := strconv.Atoi(start[i])
			pods = append(pods, podRange{stage: startType, valStart: valStart, valEnd: valStart})
		}

	} else if question == 2 {
		for i := 0; i < len(start); i += 2 {
			valStart, _ := strconv.Atoi(start[i])
			valRange, _ := strconv.Atoi(start[i+1])
			pods = append(pods, podRange{stage: startType, valStart: valStart, valEnd: valStart + valRange - 1})
		}

	}

	transformers := make(map[string]transformer)

	typeRegex := regexp.MustCompile("([A-Za-z]+)-to-([A-Za-z]+) map:")
	dataRegex := regexp.MustCompile(`(\d+) (\d+) (\d+)`)

	currentKey := ""
	var currentCharts []Chart
	var currentTransformer transformer
	for i := 2; i < len(lines); i++ {
		data := dataRegex.FindStringSubmatch(lines[i])
		if data != nil {
			rng, _ := strconv.Atoi(data[3])
			source, _ := strconv.Atoi(data[2])
			target, _ := strconv.Atoi(data[1])
			currentCharts = append(currentCharts, Chart{begin: source, end: source + rng - 1, diff: source - target})
		} else {
			if len(currentCharts) > 0 {
				currentCharts = storeCharts(currentCharts, currentTransformer, currentKey, transformers)
			}
		}
		mapping := typeRegex.FindStringSubmatch(lines[i])
		if mapping != nil {
			currentKey = mapping[1]
			currentTransformer = transformer{inputType: mapping[1], outputType: mapping[2], chart: make([]Chart, 0)}
		}
	}
	if len(currentCharts) > 0 {
		currentCharts = storeCharts(currentCharts, currentTransformer, currentKey, transformers)
	}
	return transformers, pods
}

func storeCharts(currentCharts []Chart, currTransformer transformer, key string, transformers map[string]transformer) []Chart {
	slices.SortFunc(currentCharts, func(a, b Chart) int {
		return cmp.Compare(a.begin, b.begin)
	})
	currTransformer.chart = currentCharts
	transformers[key] = currTransformer
	return []Chart{}
}
