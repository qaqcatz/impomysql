package ckstable

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.Wrap(err, "[PathExists]file stat error")
}

func execCmd(cmdStr string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	return string(out), err
}