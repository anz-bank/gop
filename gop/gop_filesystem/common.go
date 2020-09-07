package gop_filesystem

import (
	"github.com/joshcarp/gop/app"
	"github.com/spf13/afero"
)

func New(fs afero.Fs, appconfig app.AppConfig) GOP {
	return GOP{
		AppConfig: appconfig,
		fs:        fs,
	}
}
