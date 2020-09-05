package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-git/v5"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type retriever func(res *pbmod.KeyValue) (err error)

func (a server) retrieveGit(res *pbmod.KeyValue) error {
	var auth *http.BasicAuth
	store := memory.NewStorage()
	if a.username != "" {
		auth = &http.BasicAuth{
			Username: a.username,
			Password: a.token,
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
	return retrieveFie(&res.Value, reader)
}

func (a server) retrieveSyslPb(res *pbmod.KeyValue) error {
	file, err := os.Open(path.Join(a.saveLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	if err := retrieveFie(res.Extra, file); err != nil {
		return err
	}
	return a.retrieveFile(res)
}

func (a server) retrieveFile(res *pbmod.KeyValue) error {
	file, err := os.Open(path.Join(a.saveLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	res.Imported = true
	return retrieveFie(&res.Value, file)
}

func retrieveFie(res *string, file io.Reader) error {
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fin := string(contents)
	*res = fin
	return nil
}

func (a server) retrieveSyslPB(res *pbmod.KeyValue) error {
	file, err := os.Open(path.Join(a.saveLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	res.Imported = true
	return retrieveFie(res.Extra, file)
}
