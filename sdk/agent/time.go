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

func date() int {
	return getDate(currentTime())
}

func getDate(t time.Time) int {
	year, month, day := t.Date()
	return year*10000 + int(month)*100 + day
}
