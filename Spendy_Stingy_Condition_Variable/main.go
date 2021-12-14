package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock           = sync.Mutex{}
	money          = 100
	wg             = sync.WaitGroup{}
	moneyDeposited = sync.NewCond(&lock)
)

func stingy() {
	for i := 0; i < 1000; i++ {
		lock.Lock()
		money += 10
		fmt.Println("Stingy Sees Balance of", money)
		moneyDeposited.Signal()
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	fmt.Println("Stingy Done")

}

func spendy() {
	for i := 0; i < 1000; i++ {
		lock.Lock()
		for money-20 < 200 {
			moneyDeposited.Wait()
		}
		money -= 20
		fmt.Println("Spendy sees balance of", money)
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	fmt.Println("Stingy Done")

}

func main() {
	go stingy()
	go spendy()
	time.Sleep(3 * time.Second)
	fmt.Println("Money After Stingy Spendy", money)

}
