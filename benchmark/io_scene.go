package benchmark

import "benchmark/ellyn_agent"

import (
	"os"
	"path/filepath"
)

var devNull *os.File
var tmpFile *os.File

func init() {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 3, nil)
	defer ellyn_agent.Agent.Pop(_ellynCtx, nil)
	ellyn_agent.Agent.SetBlock(_ellynCtx, 0, 8)
	f, err := os.Open(os.DevNull)
	if err != nil {
		ellyn_agent.Agent.SetBlock(_ellynCtx, 1, 11)
		panic(err)
	}
	ellyn_agent.Agent.SetBlock(_ellynCtx, 2, 9)
	devNull = f
	f, err = os.OpenFile(filepath.Join(os.TempDir(), "bench.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		ellyn_agent.Agent.SetBlock(_ellynCtx, 3, 12)
		panic(err)
	}
	ellyn_agent.Agent.SetBlock(_ellynCtx, 4, 10)
	tmpFile = f
}

func Write2DevNull(content string) {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 4, nil)
	defer ellyn_agent.Agent.Pop(_ellynCtx, nil)
	ellyn_agent.Agent.SetBlock(_ellynCtx, 0, 13)
	_, _ = devNull.Write([]byte(content))
}

func Write2TmpFile(content string) {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 5, nil)
	defer ellyn_agent.Agent.Pop(_ellynCtx, nil)
	ellyn_agent.Agent.SetBlock(_ellynCtx, 0, 14)
	_, _ = tmpFile.Write([]byte(content))
}
