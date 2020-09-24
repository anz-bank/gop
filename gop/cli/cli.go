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
	Resolver  gop.Resolver
}

func New(local gop.Gopper, cache gop.Gopper, proxy gop.Retriever, github gop.Retriever, git gop.Retriever, cacheFile string, resolver gop.Resolver) Retriever {
	return Retriever{
		cacheFile: cacheFile,
		local:     local,
		cache:     cache,
		proxy:     proxy,
		github:    github,
		git:       git,
		Resolver:  resolver,
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
	gh := retriever_github.New(token)
	return New(
		gop_filesystem.New(fs, "."),
		cache,
		proxy,
		gh,
		retriever_git.New(token),
		cacheFile,
		gh.ResolveHash)
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
		orig := resource
		if resource == "github.com/joshcarp/sysl-1/sysl-1.sysl@v1.0.0" {
			println(orig)
		}
		resource, err = gop.LoadVersion(r.cache, r.Resolver, r.cacheFile, resource)
		if resource == "ee5c8cd2b97ba24226cb86556c17ad4f852915f2" {
			println(orig)
		}
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
