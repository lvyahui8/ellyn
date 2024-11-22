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

var levelBytes = [][]byte{
	[]byte("[Info]"),
	[]byte("[Warn]"),
	[]byte("[Error]"),
}

func (level logLevel) strBytes() []byte {
	return levelBytes[level]
}

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

var linePool = &sync.Pool{
	New: func() any {
		return &logLine{
			buf: make([]byte, 0, 128),
		}
	},
}

type logLine struct {
	buf []byte
}

func (l *logLine) Recycle() {
	l.buf = l.buf[:0]
	linePool.Put(l)
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
	l.logFormatMsg(Error, format, args)
}

func (l *asyncLogger) Warn(format string, args ...any) {
	if l.currentLevel > Warn {
		return
	}
	l.logFormatMsg(Warn, format, args)
}

func (l *asyncLogger) Info(format string, args ...any) {
	if l.currentLevel > Info {
		return
	}
	l.logFormatMsg(Info, format, args)
}

func (l *asyncLogger) logFormatMsg(level logLevel, format string, args []any) {
	line := linePool.Get().(*logLine)
	line.buf = append(line.buf, currDatetime...)
	line.buf = append(line.buf, ' ')
	line.buf = append(line.buf, level.strBytes()...)
	line.buf = append(line.buf, ' ')
	if len(args) == 0 {
		line.buf = append(line.buf, utils.String.String2bytes(format)...)
	} else {
		message := fmt.Sprintf(format, args...) // Sprintf性能较差，可以进一步优化为自行实现
		line.buf = append(line.buf, utils.String.String2bytes(message)...)
	}
	line.buf = append(line.buf, '\n')
	suc := l.logQueue.Enqueue(line)
	if !suc {
		line.Recycle()
	}
}

func (l *asyncLogger) outputLine(line *logLine) {
	_, _ = l.file.Write(line.buf)
	line.Recycle()
}
