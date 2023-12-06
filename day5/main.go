package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type pod struct {
	stage string
	value int
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
	transformers, pods := loadData()
	fmt.Println(transformers, pods)
	fmt.Println("finish load")
	var answer1 int
	for _, pod := range pods {
		res, err := pod.convert("location", *transformers)
		if err != nil {
			panic(err)
		}
		if res < answer1 || answer1 == 0 {
			answer1 = res
		}
	}
	fmt.Println(answer1)
}

func loadData() (*map[string]transformer, []pod) {
	file, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")

	var pods []pod
	start := strings.Split(lines[0], " ")
	startType := start[0][:len(start[0])-2]
	start = start[1:]
	for _, v := range start {
		val, _ := strconv.Atoi(v)
		pods = append(pods, pod{stage: startType, value: val})
	}

	// fmt.Println(pods)

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
		}
		mapping := typeRegex.FindStringSubmatch(lines[i])
		if mapping != nil {
			if len(currentCharts) > 0 {
				currentTransformer.chart = currentCharts
				transformers[currentKey] = currentTransformer
				currentCharts = []Chart{}
			}
			fmt.Println("new map ", mapping[0])
			currentKey = mapping[1]
			currentTransformer = transformer{inputType: mapping[1], outputType: mapping[2], chart: make([]Chart, 0)}
		}
	}
	if len(currentCharts) > 0 {
		currentTransformer.chart = currentCharts
		transformers[currentKey] = currentTransformer
		currentCharts = []Chart{}
	}
	return &transformers, pods
}

func (p *pod) convert(target string, t map[string]transformer) (int, error) {
	for p.stage != target {
		next, match := t[p.stage]
		if match == false {
			return 0, fmt.Errorf("Not found")
		}
		// fmt.Println(next)
		// fmt.Println(p)

		for _, v := range next.chart {
			if p.value >= v.begin && p.value <= v.end {
				p.value = p.value - v.diff
				break
			}
		}
		p.stage = next.outputType
	}
	return p.value, nil
}
