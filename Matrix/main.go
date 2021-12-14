package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const matrixSize = 250

const numberOfThreads int = 8

var waitGroup = sync.WaitGroup{}

var (
	matrixA = [matrixSize][matrixSize]int{}

	matrixB = [matrixSize][matrixSize]int{}

	result         = [matrixSize][matrixSize]int{}
	integerChannel = make(chan int, 1000)
)

func generateRandomMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			matrix[row][col] += rand.Intn(10) - 5
		}
	}
}

func workOutRow(rowChannel chan int) {
	for row := range rowChannel {
		for col := 0; col < matrixSize; col++ {
			for i := 0; i < matrixSize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
	}
	waitGroup.Done()

}

func main() {
	fmt.Println("Working")
	count := 0
	start := time.Now()
	for i := 0; i < 100; i++ {
		for row := 0; row < matrixSize; row++ {
			count++
			go workOutRow(integerChannel)
		}
	}

	waitGroup.Add(count)
	for i := 0; i < 100; i++ {
		generateRandomMatrix(&matrixA)
		generateRandomMatrix(&matrixB)
		for row := 0; row < matrixSize; row++ {
			integerChannel <- row

		}
	}
	close(integerChannel)
	fmt.Println(count)
	elapsed := time.Since(start)
	waitGroup.Wait()
	fmt.Println("Done")
	fmt.Println(elapsed)

}
