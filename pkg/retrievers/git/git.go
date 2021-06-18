// +build !wasm,!js

package git

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/anz-bank/gop/pkg/gop"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Retriever struct {
	// tokens is a key/value pairs of <host>, <personal access token>, e.g. { "github.com": "abcdef" }
	tokens         map[string]string
	privateKeyFile string
	password       string
}

/* New returns a retriever */
func New(tokens map[string]string, privateKeyFile string, password string) Retriever {
	if tokens == nil {
		tokens = map[string]string{}
	}
	return Retriever{tokens: tokens, privateKeyFile: privateKeyFile, password: password}
}

func (a Retriever) getToken(resource string) string {
	u, _ := url.Parse("https://" + resource)
	return a.tokens[u.Host]
}

func (a Retriever) Retrieve(resource string) ([]byte, bool, error) {
	var r *git.Repository
	store := memory.NewStorage()
	fs := memfs.New()
	repo, path, version, err := gop.ProcessRequest(resource)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.BadRequestError, err)
	}

	if b := a.getToken(resource); b != "" {
		auth := &http.BasicAuth{
			Username: "gop",
			Password: b,
		}
		url := "https://" + repo + ".git"
		r, err = git.Clone(store, fs, &git.CloneOptions{
			URL:  url,
			Auth: auth,
		})
		if err != nil {
			return nil, false, fmt.Errorf("%s, git clone %s via PAT, %w", gop.GitCloneError, url, err)
		}
	} else if a.privateKeyFile != "" {
		_, err = os.Stat(a.privateKeyFile)
		if err != nil {
			return nil, false, fmt.Errorf("read file %s failed %s\n", a.privateKeyFile, err.Error())
		}

		publicKeys, err := ssh.NewPublicKeysFromFile("git", a.privateKeyFile, a.password)
		if err != nil {
			return nil, false, fmt.Errorf("generate publickeys failed: %s\n", err.Error())
		}
		url := "ssh://" + repo + ".git"
		r, err = git.Clone(store, fs, &git.CloneOptions{
			URL:  url,
			Auth: publicKeys,
		})
		if err != nil {
			return nil, false, fmt.Errorf("%s, git clone %s via SSH key, %w", gop.GitCloneError, url, err)
		}
	} else {
		auth, err := ssh.NewSSHAgentAuth("git")
		if err != nil {
			return nil, false, fmt.Errorf("set up ssh-agent auth failed %s\n", err.Error())
		}
		url := "ssh://" + repo + ".git"
		r, err = git.Clone(store, fs, &git.CloneOptions{
			URL:  url,
			Auth: auth,
		})
		if err != nil {
			return nil, false, fmt.Errorf("%s, git clone %s via SSH agent, %w", gop.GitCloneError, url, err)
		}
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
