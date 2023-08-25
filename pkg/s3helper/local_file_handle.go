package s3helper

import (
	"errors"
	"os"
)

type LocalFileHandle struct {
	Path      string
	IsDestory bool
}

func (r *LocalFileHandle) Destory() error {
	if !isFileExist(r.Path) {
		r.IsDestory = true
		return nil
	} else {
		err := os.Remove(r.Path)
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
