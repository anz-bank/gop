package basic_gopper

import (
	"fmt"

	"github.com/joshcarp/gop/app"
	gop2 "github.com/joshcarp/gop/gen/pkg/servers/gop"
	"github.com/joshcarp/gop/gop/gop_filesystem"
	"github.com/joshcarp/gop/gop/gop_gcs"
	"github.com/joshcarp/gop/gop/retriever/retriever_git"
	"github.com/spf13/afero"
)

type Retriever struct {
	fs  *gop_filesystem.GOP
	gcs *gop_gcs.GOP
	git *retriever_git.Retriever
}

func (a Retriever) Retrieve(res *gop2.Object) error {
	if a.fs != nil {
		if err := a.fs.Retrieve(res); err != nil {
			fmt.Println(err)
		}
	}
	if a.gcs != nil {
		if err := a.gcs.Retrieve(res); err != nil {
			fmt.Println(err)
		}
	}
	if res.Content == "" {
		return a.git.Retrieve(res)
	}
	return nil
}

func NewGopper(a app.AppConfig) (*Retriever, error) {
	r := Retriever{}
	switch a.FsType {
	case "os":
		fs := gop_filesystem.New(afero.NewOsFs(), a)
		r.fs = &fs
	case "mem", "memory":
		fs := gop_filesystem.New(afero.NewMemMapFs(), a)
		r.fs = &fs
	case "gcs":
		gcs := gop_gcs.New(a)
		r.gcs = &gcs
	}

	return &r, nil
}
