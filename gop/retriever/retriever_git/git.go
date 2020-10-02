// +build !wasm,!js

package retriever_git

import (
	"fmt"
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
		return nil, false, fmt.Errorf("%s: %w", gop.BadRequestError, err)
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
		return nil, false, fmt.Errorf("%s, %w", gop.GitCloneError, err)
	}
	h, err := r.ResolveRevision(plumbing.Revision(version))
	if err != nil {
		return nil, false, fmt.Errorf("%s, %w", gop.GitCloneError, err)
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.GitCloneError, err)
	}
	if err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(h.String()),
	}); err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.GitCheckoutError, err)
	}
	commit, err := r.CommitObject(*h)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.GitCheckoutError, err)
	}
	f, err := commit.File(path)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.FileNotFoundError, err)
	}
	reader, err := f.Reader()
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.FileNotFoundError, err)
	}
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.FileNotFoundError, err)
	}
	return b, false, nil
}
