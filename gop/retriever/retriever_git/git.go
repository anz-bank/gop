package retriever_git

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gop"
)

type Retriever struct {
	AppConfig app.AppConfig
}

func New(appConfig app.AppConfig) Retriever {
	return Retriever{AppConfig: appConfig}
}

func (a Retriever) Retrieve(repo, resource, version string) (gop.Object, bool, error) {
	var auth *http.BasicAuth
	var res gop.Object
	store := memory.NewStorage()
	fs := memfs.New()
	if a.AppConfig.Username != "" {
		auth = &http.BasicAuth{
			Username: a.AppConfig.Username,
			Password: a.AppConfig.Token,
		}
	}
	r, err := git.Clone(store, fs, &git.CloneOptions{
		URL:  "https://" + repo + ".git",
		Auth: auth,
	})
	h, err := r.ResolveRevision(plumbing.Revision(version))
	if err != nil {
		return gop.Object{}, false, app.CreateError(app.CacheAccessError, "Failed to clone repository", err)
	}
	res = app.New(repo, resource, h.String())
	w, err := r.Worktree()
	if err != nil {
		return gop.Object{}, false, app.CreateError(app.CacheAccessError, "Failed to clone repository", err)
	}
	if err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(h.String()),
	}); err != nil {
		return gop.Object{}, false, app.CreateError(app.CacheReadError, "Failed to checkout version", err)
	}
	commit, err := r.CommitObject(*h)
	if err != nil {
		return gop.Object{}, false, app.CreateError(app.CacheReadError, "Failed to checkout version", err)
	}
	f, err := commit.File(resource)
	if err != nil {
		return gop.Object{}, false, app.CreateError(app.CacheReadError, "File does not exist", err)
	}
	reader, err := f.Reader()
	if err != nil {
		return gop.Object{}, false, app.CreateError(app.CacheReadError, "Error reading file", err)
	}
	if err := app.ScanIntoString(&res.Content, reader); err != nil {
		return gop.Object{}, false, app.CreateError(app.CacheReadError, "Error reading file", err)
	}
	return res, false, nil
}
