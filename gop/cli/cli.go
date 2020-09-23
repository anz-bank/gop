package cli

import (
	"fmt"

	"github.com/joshcarp/gop/gop"
	"github.com/joshcarp/gop/gop/gop_filesystem"
	"github.com/joshcarp/gop/gop/retriever/retriever_git"
	"github.com/joshcarp/gop/gop/retriever/retriever_github"
	"github.com/joshcarp/gop/gop/retriever/retriever_proxy"
	"github.com/spf13/afero"
)

/*
1. local retrieve from current project
2. cache retrieve
3. proxy retrieve -> cache
4. github retrieve -> cache
4. git retrieve -> cache
*/

/* Retriever Is a CLI retriever that can be used for retrieving and caching for cli tools that require remote imports */
type Retriever struct {
	cacheFile string
	local     gop.Retriever
	cache     gop.Gopper
	proxy     gop.Retriever
	http      gop.Retriever
	github    gop.Retriever
	git       gop.Retriever
}

func New(local gop.Gopper, cache gop.Gopper, proxy, github, git gop.Retriever, cacheFile string) Retriever {
	return Retriever{
		cacheFile: cacheFile,
		local:     local,
		cache:     cache,
		proxy:     proxy,
		github:    github,
		git:       git,
	}
}

func Default(fs afero.Fs, cacheFile, cacheDir string, proxyURL string, token map[string]string) Retriever {
	var cache gop.Gopper
	var proxy gop.Retriever
	if cacheDir != "" {
		cache = gop_filesystem.New(fs, cacheDir)
	}
	if proxyURL != "" {
		proxy = retriever_proxy.New(proxyURL)
	}
	return New(
		gop_filesystem.New(fs, "."),
		cache,
		proxy,
		retriever_github.New(token),
		retriever_git.New(token),
		cacheFile)
}

/* Retrieve implements the retriever interface */
func (r Retriever) Retrieve(resource string) ([]byte, bool, error) {
	var content []byte
	var err error
	var cummulative error
	if r.local != nil {
		content, _, err = r.local.Retrieve(resource)
		if !(err != nil || content == nil || len(content) == 0) {
			return content, false, nil
		}
		cummulative = fmt.Errorf("%s: %w\n", gop.FileNotFoundError, err)
	}
	if r.cache != nil {
		resource, err := gop.LoadVersion(r.cache, gop.ResolveHash, r.cacheFile, resource)
		if err != nil {
			return nil, false, err
		}

		content, _, err = r.cache.Retrieve(resource)
		if !(err != nil || content == nil || len(content) == 0) {
			return content, false, nil
		}

		cummulative = fmt.Errorf("%s: %w\n", cummulative, err)
		defer func() {
			if repo, _, ver, _ := gop.ProcessRequest(resource); repo != "" && ver != "" && len(content) != 0 {
				r.cache.Cache(resource, content)
			}
		}()
	}
	if r.proxy != nil {
		content, _, err = r.proxy.Retrieve(resource)
		if !(err != nil || content == nil || len(content) == 0) {
			return content, false, nil
		}
		cummulative = fmt.Errorf("%s: %w\n", cummulative, err)
	}

	if r.github != nil {
		content, _, err = r.github.Retrieve(resource)
		if !(err != nil || content == nil || len(content) == 0) {
			return content, false, nil
		}
		cummulative = fmt.Errorf("%s: %w\n", cummulative, err)
	}
	if r.git != nil {
		content, _, err = r.git.Retrieve(resource)
		if !(err != nil || content == nil || len(content) == 0) {
			return content, false, nil
		}
		cummulative = fmt.Errorf("%s: %w\n", cummulative, err)
	}
	return content, false, cummulative
}
