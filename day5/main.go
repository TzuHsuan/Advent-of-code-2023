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
	t2, p2 := loadData2()
	var answer2 int
	for len(p2) != 0 {
		for p2[0].stage != "location" {
			currentCharts := t2[p2[0].stage].chart
			fmt.Println(p2[0])
			for cp, pp := 0, p2[0].valStart; pp < p2[0].valEnd; {
				cc := currentCharts[cp]
				if pp > cc.end {
					cp += 1
					if cp == len(currentCharts) {
						// fmt.Println("skipping ", p2[0])
						p2[0].stage = t2[p2[0].stage].outputType
						break
					}
					continue
				}

				np := podRange{}

				//Starts outside of map, No change is value
				if pp < cc.begin {
					np.valStart = pp
					//Entire range is outside map
					if p2[0].valEnd < cc.begin {
						fmt.Println("before map", p2[0])
						np.valEnd = p2[0].valEnd
						np.stage = t2[p2[0].stage].outputType
						p2[0] = np
						break
					} else { //Extends into map, create split
						np = podRange{valStart: pp, valEnd: cc.begin - 1, stage: t2[p2[0].stage].outputType}
						p2 = append(p2, np)
						p2[0].stage = t2[p2[0].stage].outputType
						p2[0].valStart = cc.begin
						pp = cc.begin
					}
				} else { //Starts in map
					//Entire range is in map
					if p2[0].valEnd <= cc.end {
						np = podRange{valStart: pp - cc.diff, valEnd: p2[0].valEnd - cc.diff, stage: t2[p2[0].stage].outputType}
						p2[0] = np
						break
					} else { //Extends out of map, create split
						np = podRange{valStart: pp - cc.diff, valEnd: cc.end - cc.diff, stage: t2[p2[0].stage].outputType}
						p2 = append(p2, np)
						p2[0].valStart = cc.end + 1
						pp = cc.end + 1
					}
				}
			}
		}
		if p2[0].valStart < answer2 || answer2 == 0 {
			answer2 = p2[0].valStart
		}
		fmt.Println(p2[0])
		p2 = p2[1:]
	}
	fmt.Println("Answer 2 ", answer2)
}

func loadData2() (map[string]transformer, []podRange) {
	file, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")

	var pods []podRange
	start := strings.Split(lines[0], " ")
	startType := start[0][:len(start[0])-2]
	start = start[1:]
	for i := 0; i < len(start); i += 2 {
		valStart, _ := strconv.Atoi(start[i])
		valRange, _ := strconv.Atoi(start[i+1])
		pods = append(pods, podRange{stage: startType, valStart: valStart, valEnd: valStart + valRange - 1})
	}

	fmt.Println(pods)

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
				slices.SortFunc(currentCharts, func(a, b Chart) int {
					return cmp.Compare(a.begin, b.begin)
				})
				currentTransformer.chart = currentCharts
				transformers[currentKey] = currentTransformer
				currentCharts = []Chart{}
			}
		}
		mapping := typeRegex.FindStringSubmatch(lines[i])
		if mapping != nil {
			fmt.Println("new map ", mapping[0])
			currentKey = mapping[1]
			currentTransformer = transformer{inputType: mapping[1], outputType: mapping[2], chart: make([]Chart, 0)}
		}
	}
	if len(currentCharts) > 0 {
		slices.SortFunc(currentCharts, func(a, b Chart) int {
			return cmp.Compare(a.begin, b.begin)
		})
		currentTransformer.chart = currentCharts
		transformers[currentKey] = currentTransformer
		currentCharts = []Chart{}
	}
	// fmt.Println(transformers)
	return transformers, pods
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
