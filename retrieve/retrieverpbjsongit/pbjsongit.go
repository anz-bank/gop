package retrieverpbjsongit

import (
	"fmt"

	"github.com/joshcarp/pb-mod/retrieve/retrieverfs"
	"github.com/joshcarp/pb-mod/retrieve/retrievergit"

	"github.com/joshcarp/pb-mod/config"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type RetrieveFilePBJsonGit struct {
	AppConfig config.AppConfig
	fs        retrieverfs.RetrieverFstruct
	git       retrievergit.RetrieverGit
}

func (a RetrieveFilePBJsonGit) Retrieve(res *pbmod.Object) error {
	if err := a.fs.Retrieve(res); err != nil {
		fmt.Println(err)
	}
	if res.Value == "" {
		return a.git.Retrieve(res)
	}
	return nil
}
