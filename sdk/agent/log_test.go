package agent

import "testing"

func TestLogger_Info(t *testing.T) {
	initLogger()
	for i := 0; i < 100000; i++ {
		log.Info("hello world")
	}
	log.file.w.Flush()
}
