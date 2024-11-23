package agent

import (
	"bufio"
	"fmt"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// LogFileMaxSize 日志文件最大100M
	LogFileMaxSize     = 100 * 1024 * 1024
	LogFileMaintainDay = 7
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
	size int64
	date int
	w    *bufio.Writer
	file *os.File
}

func newRotateFile() *logfile {
	f := &logfile{}
	f.checkRotate()
	f.autoClean()
	return f
}

func (f *logfile) autoClean() {
	go func() {
		//ticker := time.NewTicker(24 * time.Hour)
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case cur := <-ticker.C:
				minDay := getDate(cur.Add(-LogFileMaintainDay * 24 * time.Hour))
				baseFile := f.getBaseLogFile()
				dir := filepath.Dir(baseFile)
				items, err := os.ReadDir(dir)
				if err != nil {
					break
				}
				for _, item := range items {
					itemName := filepath.Join(dir, item.Name())
					if !strings.HasPrefix(itemName, baseFile) {
						continue
					}
					begin := len(baseFile)
					if len(itemName) < begin+10 {
						continue
					}
					dayStr := itemName[begin+1 : begin+8+1]
					day, err := strconv.Atoi(dayStr)
					if err != nil {
						continue
					}
					if day < minDay {
						// 可以清理
						utils.OS.Remove(itemName)
					}
				}
			}
		}
	}()
}

func (f *logfile) Write(p []byte) (n int, err error) {
	n, err = f.w.Write(p)
	f.size += int64(n)

	f.checkRotate()
	return
}

func (f *logfile) checkRotate() {
	if f.date != date() {
		f.date = date()
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
		err = os.Rename(f.file.Name(), f.dumpFileName())
		if err != nil {
			return
		}
	}
	// 指向新文件
	logFile := f.getBaseLogFile()
	initSize := int64(0)
	fileInfo, statErr := os.Stat(logFile)
	if statErr == nil {
		initSize = fileInfo.Size()
	}
	file, err := os.OpenFile(logFile,
		os.O_CREATE|os.O_APPEND|os.O_RDWR, 644)
	asserts.IsNil(err)
	w := bufio.NewWriter(file)
	f.w = w
	f.file = file
	f.size = initSize
	return
}

func (f *logfile) dumpFileName() string {
	file := f.getBaseLogFile()
	dir := filepath.Dir(file)
	items, err := os.ReadDir(dir)
	asserts.IsNil(err)
	prefix := fmt.Sprintf("%s.%d.", file, f.date)
	maxIdx := -1
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		itemName := filepath.Join(dir, item.Name())
		if strings.HasPrefix(itemName, prefix) {
			idx, parseErr := strconv.Atoi(itemName[len(prefix):])
			if parseErr != nil {
				continue
			}
			if idx > maxIdx {
				maxIdx = idx
			}
		}
	}
	return fmt.Sprintf("%s%d", prefix, maxIdx+1)
}

func (f *logfile) getBaseLogFile() string {
	wd, err := os.Getwd()
	asserts.IsNil(err)
	logPath := filepath.Join(wd, "logs")
	if utils.OS.NotExists(logPath) {
		utils.OS.MkDirs(logPath)
	}
	return filepath.Join(logPath, "run.log")
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

type KVItem interface {
	formatLogBytes() []byte
}

var schemaPool = &sync.Pool{New: func() any {
	return &Schema{
		res: make([]byte, 0, 64),
	}
}}

type Schema struct {
	items []KVItem
	res   []byte
}

func (s *Schema) Build() []byte {
	for i, item := range s.items {
		if i != 0 {
			s.res = append(s.res, '|')
		}
		s.res = append(s.res, item.formatLogBytes()...)
	}
	return s.res
}

func (s *Schema) Recycle() {
	s.items = s.items[:0]
	s.res = s.res[:0]
	schemaPool.Put(s)
}

type KV[VT int | string | bool] struct {
	key string
	val VT
}

func (item *KV[VT]) formatLogBytes() (res []byte) {
	res = append(res, item.key...)
	res = append(res, ':')
	switch n := any(item.val).(type) {
	case int:
		res = append(res, strconv.Itoa(n)...)
	case string:
		res = append(res, n...)
	case bool:
		if n {
			res = append(res, 'Y')
		} else {
			res = append(res, 'N')
		}
	}
	return
}

func V[VT int | string | bool](key string, val VT) *KV[VT] {
	return &KV[VT]{key: key, val: val}
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
		last := currTime.Unix()
		for {
			line, suc := l.logQueue.Dequeue()
			if !suc {
				time.Sleep(time.Microsecond)
				now := currTime.Unix()
				if now-last > 1 {
					// 超过1s无日志输出
					l.flush()
				}
				continue
			}
			last = currTime.Unix()
			l.outputLine(line)
		}
	}()
}

func (l *asyncLogger) InfoKV(items ...KVItem) {
	line := l.newLine(Info)
	s := schemaPool.Get().(*Schema)
	s.items = items
	line.buf = append(line.buf, s.Build()...)
	l.emitLine(line)
	s.Recycle()
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
	line := l.newLine(level)
	if len(args) == 0 {
		line.buf = append(line.buf, utils.String.String2bytes(format)...)
	} else {
		message := fmt.Sprintf(format, args...) // Sprintf性能较差，可以进一步优化为自行实现
		line.buf = append(line.buf, utils.String.String2bytes(message)...)
	}
	l.emitLine(line)
}

func (l *asyncLogger) emitLine(line *logLine) {
	line.buf = append(line.buf, '\n')
	suc := l.logQueue.Enqueue(line)
	if !suc {
		line.Recycle()
	}
}

func (l *asyncLogger) newLine(level logLevel) *logLine {
	line := linePool.Get().(*logLine)
	line.buf = append(line.buf, currDatetime...)
	line.buf = append(line.buf, ' ')
	line.buf = append(line.buf, level.strBytes()...)
	line.buf = append(line.buf, ' ')
	return line
}

func (l *asyncLogger) outputLine(line *logLine) {
	_, _ = l.file.Write(line.buf)
	line.Recycle()
}

func (l *asyncLogger) flush() {
	_ = l.file.w.Flush()
}
