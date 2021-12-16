package main

import (
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	totalAccounts  = 50000
	maxAmountMoved = 10
	initialMoney   = 100
	threads        = 4
)

func perform_movements(ledger *[totalAccounts]int32, locks *[totalAccounts]sync.Locker, totalTransactions *int64) {
	for {
		accountA := rand.Intn(totalAccounts)
		accountB := rand.Intn(totalAccounts)
		for accountA == accountB {
			accountB = rand.Intn(totalAccounts)
		}
		amountToMove := rand.Int31n(maxAmountMoved)
		toLock := []int{accountA, accountB}
		sort.Ints(toLock)

		locks[toLock[0]].Lock()
		locks[toLock[1]].Lock()
		atomic.AddInt32(&ledger[accountA], -amountToMove)
		atomic.AddInt32(&ledger[accountB], amountToMove)
		atomic.AddInt64(totalTransactions, 1)
		locks[toLock[0]].Unlock()
		locks[toLock[1]].Unlock()
	}

}

func main() {
	println("Total Accounts :", totalAccounts, "total threads", threads, "using spinLocks")
	var ledger [totalAccounts]int32
	var locks [totalAccounts]sync.Locker
	var totalTrans int64
	for i := 0; i < totalAccounts; i++ {
		ledger[i] = initialMoney
		locks[i] = NewSpinLock()
	}

	for i := 0; i < threads; i++ {
		go perform_movements(&ledger, &locks, &totalTrans)
	}
	for {
		time.Sleep(2000 * time.Millisecond)
		var sum int32
		for i := 0; i < totalAccounts; i++ {
			locks[i].Lock()
		}
		for i := 0; i < totalAccounts; i++ {
			sum += ledger[i]
		}
		for i := 0; i < totalAccounts; i++ {
			locks[i].Unlock()
		}
		println(totalTrans, sum)
	}
}
