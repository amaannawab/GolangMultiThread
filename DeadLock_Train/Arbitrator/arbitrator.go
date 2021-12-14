package arbitrator

import (
	"DeadLock_Train/common"
	"sync"
	"time"
)

var (
	controller = sync.Mutex{}
	cond       = sync.NewCond(&controller)
)

func allFree(intersectionsToLock []*common.Intersection) bool {
	for _, it := range intersectionsToLock {
		if it.LockedBy >= 0 {
			return false
		}
	}
	return true
}

func lockIntersectionsInDistance(id, reserveStart, reserveEnd int, crossings []*common.Crossing) {
	intersectionsToLock := []*common.Intersection{}
	for _, crossing := range crossings {
		if reserveEnd >= crossing.Position && reserveStart <= crossing.Position && crossing.Intersection.LockedBy != id {
			intersectionsToLock = append(intersectionsToLock, crossing.Intersection)
		}
	}

	controller.Lock()
	for !allFree(intersectionsToLock) {
		cond.Wait()
	}

	for _, intersection := range intersectionsToLock {

		intersection.LockedBy = id
		time.Sleep(10 * time.Millisecond)
	}
	controller.Unlock()

}

func MoveTrain(train *common.Train, distance int, crossings []*common.Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				lockIntersectionsInDistance(train.Id, crossing.Position, crossing.Position+train.TrainLength, crossings)
				// crossing.Intersection.Mutex.Lock()
				// crossing.Intersection.LockedBy = train.Id
			}
			back := train.Front - train.TrainLength
			if back == crossing.Position {
				controller.Lock()
				crossing.Intersection.LockedBy = -1
				controller.Unlock()
				cond.Broadcast()
			}

		}
		time.Sleep(30 * time.Millisecond)
	}
}
