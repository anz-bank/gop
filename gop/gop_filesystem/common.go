package gop_filesystem

import (
	"github.com/spf13/afero"
)

func New(fs afero.Fs, dir string) GOP {
	return GOP{
		fs:  fs,
		dir: dir,
	}
}
