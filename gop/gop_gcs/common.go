package gop_gcs

import (
	"cloud.google.com/go/storage"
	"github.com/joshcarp/gop/app"
)

type GOP struct {
	AppConfig  app.AppConfig
	upload     uploader
	downloader downloader
	client     *storage.Client
}

func New(appconfig app.AppConfig) GOP {
	return GOP{
		AppConfig:  appconfig,
		upload:     UploadFile,
		downloader: download,
	}
}
