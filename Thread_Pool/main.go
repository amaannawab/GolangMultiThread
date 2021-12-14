package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Point2D struct {
	x int
	y int
}

const numberOfThreads int = 8

var (
	r         = regexp.MustCompile(`\((\d*),(\d*)\)`)
	waitGroup = sync.WaitGroup{}
)

func findArea(inputChannel chan string) {
	for pointsStr := range inputChannel {
		var points []Point2D
		for _, p := range r.FindAllStringSubmatch(pointsStr, -1) {
			x, _ := strconv.Atoi(p[1])
			y, _ := strconv.Atoi(p[2])
			points = append(points, Point2D{x: x, y: y})
		}

		area := 0.0
		for i := 0; i < len(points); i++ {
			// shoelace algorithm is being used to calculate area of a polygon
			a, b := points[i], points[(i+1)%len(points)]
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}
		fmt.Println(math.Abs(area) / 2.0)
	}
	waitGroup.Done()
}

func main() {
	absPath, _ := filepath.Abs("./")
	count := 0

	dat, _ := ioutil.ReadFile(filepath.Join(absPath, "polygons.txt"))
	text := string(dat)
	inputChannel := make(chan string, 10000)
	for i := 0; i < numberOfThreads; i++ {
		go findArea(inputChannel)
	}
	waitGroup.Add(numberOfThreads)
	start := time.Now()
	for _, line := range strings.Split(text, "\n") {
		count++

		inputChannel <- line
	}
	close(inputChannel)
	waitGroup.Wait()
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	fmt.Println(count)

}
