package timer

import "time"

type Timer struct {
	ticker *time.Ticker
}

func StartTimer() {
	newTimer := time.NewTimer(5 * time.Second)
	currentTime := <-newTimer.C
}
