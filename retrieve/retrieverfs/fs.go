package retrieverfs

import (
	"fmt"
	"os"
	"path"

	"github.com/joshcarp/pb-mod/config"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type RetrieverFstruct struct {
	AppConfig config.AppConfig
}

func (a RetrieverFstruct) Retrieve(res *pbmod.Object) error {
	file, err := os.Open(path.Join(a.AppConfig.SaveLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	if err := config.ScanIntoString(res.Extra, file); err != nil {
		return err
	}
	return a.RetrieverFile(res)
}

func (a RetrieverFstruct) RetrieverFile(res *pbmod.Object) error {
	file, err := os.Open(path.Join(a.AppConfig.SaveLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	res.Imported = true
	return config.ScanIntoString(&res.Value, file)
}
