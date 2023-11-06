package file

import (
	"errors"
	"os"
	"path"
)

type LocalFileHandle struct {
	Path      string
	IsDestory bool
}

var (
	osRemove  = os.Remove
	osMakeDir = os.MkdirAll
	osCreate  = os.Create
)

func (r *LocalFileHandle) Destory() error {
	if !IsFileExist(r.Path) {
		r.IsDestory = true
		return nil
	} else {
		err := osRemove(r.Path)
		if err == nil {
			r.IsDestory = true
			return nil
		}
		if errors.Is(err, os.ErrNotExist) {
			r.IsDestory = true
			return nil
		}
		return err
	}
}

func IsFileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func CreateFileAndReturnFile(fullPath string) (*os.File, error) {
	err := osMakeDir(path.Dir(fullPath), 0770)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
	}
	file, err := osCreate(fullPath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
