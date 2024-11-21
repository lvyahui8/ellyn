package agent

import (
	"bufio"
	"fmt"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	LogFileMaxSize = 1 * 1024 * 1024
)

// log
// log写盘逻辑：定时、大小溢出、日期切换
var log *asyncLogger

var logInitOnce = sync.Once{}

type logLevel int

const (
	Info logLevel = iota
	Warn
	Error
)

type logfile struct {
	size int
	date int
	w    *bufio.Writer
	file *os.File
	idx  int
}

func newRotateFile() *logfile {
	f := &logfile{}
	f.checkRotate()
	return f
}

func (f *logfile) Write(p []byte) (n int, err error) {
	n, err = f.w.Write(p)
	f.size += n

	f.checkRotate()
	return
}

func (f *logfile) checkRotate() {
	if f.date != date() {
		f.date = date()
		f.idx = 0
		_ = f.rotate()
	} else if f.size >= LogFileMaxSize {
		_ = f.rotate()
	}
}

// 将当前文件重命名，创建新文件
func (f *logfile) rotate() (err error) {
	// 刷新当前文件
	if f.w != nil {
		err = f.w.Flush()
		if err != nil {
			return
		}
	}
	// 重命名文件
	if f.file != nil {
		err = f.file.Close()
		if err != nil {
			return
		}
		err = os.Rename(f.file.Name(), fmt.Sprintf("%s.%d.%d", f.file.Name(), f.date, f.idx))
		if err != nil {
			return
		}
	}
	// 指向新文件
	wd, err := os.Getwd()
	asserts.IsNil(err)
	logPath := filepath.Join(wd, "logs")
	utils.OS.MkDirs(logPath)
	file, err := os.OpenFile(filepath.Join(logPath, "run.log"),
		os.O_CREATE|os.O_APPEND|os.O_RDWR, 644)
	asserts.IsNil(err)
	w := bufio.NewWriter(file)
	f.w = w
	f.file = file
	f.size = 0
	return
}

type logLine struct {
	buf []byte
}

type asyncLogger struct {
	logQueue     *collections.RingBuffer[*logLine]
	file         *logfile
	currentLevel logLevel
}

func initLogger() {
	logInitOnce.Do(func() {
		log = &asyncLogger{
			file:     newRotateFile(),
			logQueue: collections.NewRingBuffer[*logLine](4096),
		}
		log.start()
	})
}

func (l *asyncLogger) start() {
	go func() {
		for {
			line, suc := l.logQueue.Dequeue()
			if !suc {
				time.Sleep(time.Microsecond)
				continue
			}
			l.outputLine(line)
		}
	}()
}

func (l *asyncLogger) Error(format string, args ...any) {
	if l.currentLevel > Error {
		return
	}
	line := fmt.Sprintf("[Error] "+format+"\n", args...)
	l.logQueue.Enqueue(&logLine{buf: utils.String.String2bytes(line)})
}

func (l *asyncLogger) Warn(format string, args ...any) {
	if l.currentLevel > Warn {
		return
	}
	line := fmt.Sprintf("[Warn] "+format+"\n", args...)
	l.logQueue.Enqueue(&logLine{buf: utils.String.String2bytes(line)})
}

func (l *asyncLogger) Info(format string, args ...any) {
	if l.currentLevel > Info {
		return
	}
	line := fmt.Sprintf("[Info] "+format+"\n", args...)
	l.logQueue.Enqueue(&logLine{buf: utils.String.String2bytes(line)})
}

func (l *asyncLogger) outputLine(line *logLine) {
	_, _ = l.file.Write(line.buf)
}
