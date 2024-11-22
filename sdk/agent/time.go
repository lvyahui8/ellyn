package agent

import "time"

var (
	currTime     *time.Time
	currDatetime []byte
)

func currentTime() time.Time {
	return *currTime
}

func init() {
	first := time.Now()
	refreshTime(&first)
	initTimeClock()
}

func initTimeClock() {
	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)
		for {
			select {
			case cur := <-ticker.C:
				refreshTime(&cur)
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

func refreshTime(t *time.Time) {
	currTime = t
	currDatetime = []byte(t.Format("2006-01-02 15:04:05.000"))
}
