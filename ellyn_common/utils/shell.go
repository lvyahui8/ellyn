package utils

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"os"
	"os/exec"
)

var Shell *shellUtils = &shellUtils{}

type shellUtils struct {
}

func (s *shellUtils) Exec(dir, name string, args ...string) string {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stderr = os.Stderr
	outBytes, err := cmd.Output()
	asserts.IsNil(err)
	return string(outBytes)
}
