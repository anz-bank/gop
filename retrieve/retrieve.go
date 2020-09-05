package retrieve

import (
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joshcarp/pb-mod/config"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type RetrieveFilePBJsonGit struct {
	AppConfig config.AppConfig
}

type Retriever interface {
	Retriever(res *pbmod.Module) (err error)
}

func (a RetrieveFilePBJsonGit) Retriever(res *pbmod.Module) error {
	if err := a.RetrieverFile(res); err != nil {
		fmt.Println(err)
	}
	if err := a.RetrieverFilePbJson(res); err != nil {
		fmt.Println(err)
	}
	if res.Value == "" {
		return a.RetrieverGit(res)
	}
	return nil
}

func (a RetrieveFilePBJsonGit) RetrieverGit(res *pbmod.Module) error {
	var auth *http.BasicAuth
	store := memory.NewStorage()
	if a.AppConfig.Username != "" {
		auth = &http.BasicAuth{
			Username: a.AppConfig.Username,
			Password: a.AppConfig.Token,
		}
	}
	r, err := git.Clone(store, nil, &git.CloneOptions{
		URL:   "https://" + res.Repo + ".git",
		Depth: 1,
		Auth:  auth,
	})
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(plumbing.NewHash(res.Version))
	if err != nil {
		return err
	}
	f, err := commit.File(res.Resource)
	if err != nil {
		return err
	}
	reader, err := f.Reader()
	if err != nil {
		return err
	}
	return config.ScanIntoString(&res.Value, reader)
}

func (a RetrieveFilePBJsonGit) RetrieverFilePbJson(res *pbmod.Module) error {
	file, err := os.Open(path.Join(a.AppConfig.SaveLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	if err := config.ScanIntoString(res.Extra, file); err != nil {
		return err
	}
	return a.RetrieverFile(res)
}

func (a RetrieveFilePBJsonGit) RetrieverFile(res *pbmod.Module) error {
	file, err := os.Open(path.Join(a.AppConfig.SaveLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	res.Imported = true
	return config.ScanIntoString(&res.Value, file)
}
