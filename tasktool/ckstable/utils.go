package ckstable

import (
	"github.com/pkg/errors"
	"os"
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