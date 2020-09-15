package retriever_git

import (
	"io/ioutil"

	"github.com/joshcarp/gop/gop"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Retriever struct {
	token    string
	username string
}

func New(token string, username string) Retriever {
	return Retriever{token: token, username: username}
}

func (a Retriever) Retrieve(resource string) ([]byte, bool, error) {
	var auth *http.BasicAuth
	store := memory.NewStorage()
	fs := memfs.New()
	repo, resource, version, err := gop.ProcessRequest(resource)
	if err != nil {
		return nil, false, gop.CreateError(gop.BadRequestError, "BadRequestError")
	}
	if a.token != "" {
		auth = &http.BasicAuth{
			Username: a.username,
			Password: a.token,
		}
	}
	r, err := git.Clone(store, fs, &git.CloneOptions{
		URL:  "https://" + repo + ".git",
		Auth: auth,
	})
	if err != nil {
		return nil, false, gop.CreateError(gop.CacheAccessError, "Failed to clone repository", err)
	}
	h, err := r.ResolveRevision(plumbing.Revision(version))
	if err != nil {
		return nil, false, gop.CreateError(gop.CacheAccessError, "Failed to clone repository", err)
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, false, gop.CreateError(gop.CacheAccessError, "Failed to clone repository", err)
	}
	if err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(h.String()),
	}); err != nil {
		return nil, false, gop.CreateError(gop.CacheReadError, "Failed to checkout version", err)
	}
	commit, err := r.CommitObject(*h)
	if err != nil {
		return nil, false, gop.CreateError(gop.CacheReadError, "Failed to checkout version", err)
	}
	f, err := commit.File(resource)
	if err != nil {
		return nil, false, gop.CreateError(gop.CacheReadError, "File does not exist", err)
	}
	reader, err := f.Reader()
	if err != nil {
		return nil, false, gop.CreateError(gop.CacheReadError, "Error reading file", err)
	}
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, false, gop.CreateError(gop.CacheReadError, "Error reading file", err)
	}
	return b, false, nil
}
