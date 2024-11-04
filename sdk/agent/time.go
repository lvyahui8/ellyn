package agent

import "time"

var (
	currTime *time.Time
)

func currentTime() time.Time {
	return *currTime
}

func init() {
	first := time.Now()
	currTime = &first
	initTimeClock()
}

func initTimeClock() {
	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)
		for {
			select {
			case cur := <-ticker.C:
				currTime = &cur
			}
		}
	}()
}
