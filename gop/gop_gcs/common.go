package gop_gcs

import (
	"github.com/joshcarp/gop/app"
)

type GOP struct {
	AppConfig  app.AppConfig
	upload     uploader
	downloader downloader
}

func New(appconfig app.AppConfig) GOP {
	return GOP{
		AppConfig:  appconfig,
		upload:     UploadFile,
		downloader: download,
	}
}
