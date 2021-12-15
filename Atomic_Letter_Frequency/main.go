package main

// Here we are fetching response from 100's of web pages and calculating the Frequency of Alphabets in each
// web page & counting them  with help of multiple go routines
// but in function countLettersWithAtomic we are doing it with atomic variables and
// in function countLettersWithMutex , we are using Mutex..
// The Output should be as follows,
// time took by Mutex : 31.31861025s
// time took by Atomic : 9.826484833s
// Done
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

var lock = sync.Mutex{}

func countLettersWithAtomic(url string, frequency *[26]int32, waitGroup *sync.WaitGroup) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	for i := 0; i < 20; i++ {
		for _, b := range body {
			character := strings.ToLower(string(b))

			index := strings.Index(allLetters, character)
			if index >= 0 {
				atomic.AddInt32(&frequency[index], 1)
				// frequency[index] += 1
			}

		}
	}

	// fmt.Println("Got page done")

	waitGroup.Done()
}

func countLettersWithMutex(url string, frequency *[26]int32, waitGroup *sync.WaitGroup) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	for i := 0; i < 20; i++ {
		for _, b := range body {
			character := strings.ToLower(string(b))
			lock.Lock()
			index := strings.Index(allLetters, character)
			if index >= 0 {
				frequency[index] += 1
			}
			lock.Unlock()
		}
	}

	// fmt.Println("Got page done")

	waitGroup.Done()
}

func main() {
	frequency := [26]int32{}
	frequency1 := [26]int32{}
	wg := sync.WaitGroup{}
	wg1 := sync.WaitGroup{}
	startWithMutex := time.Now()
	for i := 1000; i < 1200; i++ {
		wg.Add(1)
		go countLettersWithMutex(fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", i), &frequency, &wg)
	}
	wg.Wait()
	elapsedWithMutex := time.Since(startWithMutex)
	fmt.Println("time took by Mutex :", elapsedWithMutex)
	startWithAtomic := time.Now()
	for i := 1000; i < 1200; i++ {
		wg1.Add(1)
		go countLettersWithAtomic(fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", i), &frequency1, &wg1)
	}
	wg1.Wait()
	elapsedWithAtomic := time.Since(startWithAtomic)
	fmt.Println("time took by Atomic :", elapsedWithAtomic)
	fmt.Println("Done")

	// for i, f := range frequency {
	// 	fmt.Printf("%s ->  %d\n", string(allLetters[i]), f)
	// }
}
