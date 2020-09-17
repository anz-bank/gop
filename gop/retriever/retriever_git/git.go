package retriever_git

import (
	"io/ioutil"
	"net/url"

	"github.com/joshcarp/gop/gop"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Retriever struct {
	token map[string]string
}

/* New returns a retriever with a key/value pairs of <host>, <token> eg: New("github.com", "abcdef") */
func New(tokens map[string]string) Retriever {
	if tokens == nil {
		tokens = map[string]string{}
	}
	return Retriever{token: tokens}
}

func getToken(token map[string]string, resource string) string {
	u, _ := url.Parse("https://" + resource)
	return token[u.Host]
}

func (a Retriever) Retrieve(resource string) ([]byte, bool, error) {
	var auth *http.BasicAuth
	store := memory.NewStorage()
	fs := memfs.New()
	repo, path, version, err := gop.ProcessRequest(resource)
	if err != nil {
		return nil, false, gop.CreateError(gop.BadRequestError, "BadRequestError")
	}
	if b := getToken(a.token, resource); b != "" {
		auth = &http.BasicAuth{
			Username: "gop",
			Password: b,
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
	f, err := commit.File(path)
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
