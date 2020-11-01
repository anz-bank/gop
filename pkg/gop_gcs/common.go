package gop_gcs

type GOP struct {
	upload     uploader
	downloader downloader
	bucket     string
}

func New(bucket string) GOP {
	return GOP{
		bucket:     bucket,
		upload:     UploadFile,
		downloader: download,
	}
}
