package retriever_pbjsongit

import (
	"fmt"

	retriever_fs "github.com/joshcarp/pb-mod/gop/retriever/retriever-fs"
	retriever_git "github.com/joshcarp/pb-mod/gop/retriever/retriever-git"

	"github.com/joshcarp/pb-mod/config"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type Retriever struct {
	AppConfig config.AppConfig
	fs        retriever_fs.Retriever
	git       retriever_git.Retriever
}

func (a Retriever) Retrieve(res *pbmod.Object) error {
	if err := a.fs.Retrieve(res); err != nil {
		fmt.Println(err)
	}
	if res.Content == "" {
		return a.git.Retrieve(res)
	}
	return nil
}
