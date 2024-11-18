package agent

import (
	"bufio"
	"fmt"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"os"
	"path/filepath"
	"sync"
)

const (
	LogFileMaxSize = 1 * 1024 * 1024
)

// log
// log写盘逻辑：定时、大小溢出、日期切换
var log *logger

var logInitOnce = sync.Once{}

type logfile struct {
	sync.Mutex
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
	f.Lock()
	defer f.Unlock()

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

type logger struct {
	file *logfile
}

func initLogger() {
	logInitOnce.Do(func() {
		log = &logger{
			file: newRotateFile(),
		}
	})
}

func (l *logger) Error(format string, args ...any) {
	line := fmt.Sprintf("[Error] "+format+"\n", args...)
	_, _ = l.file.Write(utils.String.String2bytes(line))
}

func (l *logger) Info(format string, args ...any) {
	line := fmt.Sprintf("[Info] "+format+"\n", args...)
	_, _ = l.file.Write(utils.String.String2bytes(line))
}
