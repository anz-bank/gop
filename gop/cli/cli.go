package cli

import (
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

/* CLIRetriever Is a CLI retriever that can be used for retrieving and caching for cli tools that require remote imports */
type CLIRetriever struct {
	local  gop.Retriever
	cache  gop.Gopper
	proxy  gop.Retriever
	github gop.Retriever
	git    gop.Retriever
}

func New(local gop.Gopper, cache gop.Gopper, proxy, github, git gop.Retriever) CLIRetriever {
	return CLIRetriever{
		local:  local,
		cache:  cache,
		proxy:  proxy,
		github: github,
		git:    git,
	}
}

func Default(fs afero.Fs, cacheDir string, proxyURL string, token map[string]string) CLIRetriever {
	return New(
		gop_filesystem.New(fs, ""),
		gop_filesystem.New(fs, cacheDir),
		retriever_proxy.New(proxyURL),
		retriever_github.New(token),
		retriever_git.New(token))
}

/* Retrieve implements the retriever interface */
func (r CLIRetriever) Retrieve(resource string) ([]byte, bool, error) {
	var content []byte
	var err error

	content, _, err = r.local.Retrieve(resource)
	if !(err != nil || content == nil || len(content) == 0) {
		return content, false, nil
	}
	content, _, err = r.cache.Retrieve(resource)
	if !(err != nil || content == nil || len(content) == 0) {
		return content, false, nil
	}

	defer func() {
		r.cache.Cache(resource, content)
	}()

	content, _, err = r.proxy.Retrieve(resource)
	if !(err != nil || content == nil || len(content) == 0) {
		return content, false, nil
	}
	content, _, err = r.git.Retrieve(resource)
	if !(err != nil || content == nil || len(content) == 0) {
		return content, false, nil
	}
	content, _, err = r.github.Retrieve(resource)
	if !(err != nil || content == nil || len(content) == 0) {
		return content, false, nil
	}
	return content, false, err
}
