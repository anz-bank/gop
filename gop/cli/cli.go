package cli

import (
	"fmt"
	"path"

	"github.com/joshcarp/gop/gop/modules"

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
	github    gop.Retriever
	git       gop.Retriever
	versioner gop.Updater
	log       gop.Logger
}

func New(local gop.Gopper, cache gop.Gopper, proxy gop.Retriever, github gop.Retriever, versioner gop.Updater, git gop.Retriever, cacheFile string, log gop.Logger) Retriever {
	return Retriever{
		cacheFile: cacheFile,
		local:     local,
		cache:     cache,
		proxy:     proxy,
		github:    github,
		git:       git,
		versioner: versioner,
		log:       log,
	}
}

func Moduler(fs afero.Fs, cacheFile, cacheDir string, proxyURL string, token map[string]string, logger gop.Logger) Retriever {
	var cache gop.Gopper
	var proxy gop.Retriever
	if cacheDir != "" {
		cache = gop_filesystem.New(fs, cacheDir)
	}
	if proxyURL != "" {
		proxy = retriever_proxy.New(proxyURL)
	}
	gh := retriever_github.New(token)
	local := gop_filesystem.New(fs, ".")
	absModuler := path.Join(cacheDir, cacheFile)
	versioner := modules.NewLoader(local, gh.Resolve, absModuler, logger)
	return Retriever{
		cacheFile: cacheFile,
		local:     local,
		cache:     cache,
		proxy:     proxy,
		github:    modules.New(gh, absModuler),
		git:       modules.New(retriever_git.New(token), absModuler),
		versioner: versioner,
		log:       logger,
	}
}

func Default(fs afero.Fs, cacheDir string, proxyURL string, token map[string]string) Retriever {
	var cache gop.Gopper
	var proxy gop.Retriever
	var cacheFile string
	if cacheDir != "" {
		cacheFile = cacheDir + ".yaml"
		cache = gop_filesystem.New(fs, cacheDir)
	}
	if proxyURL != "" {
		proxy = retriever_proxy.New(proxyURL)
	}
	absModuler := path.Join(cacheDir, cacheFile)
	gh := retriever_github.New(token)
	return Retriever{
		cacheFile: cacheFile,
		local:     gop_filesystem.New(fs, "."),
		cache:     cache,
		proxy:     proxy,
		github:    modules.New(gh, absModuler),
		git:       modules.New(retriever_git.New(token), absModuler),
	}
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
	if r.versioner != nil {
		resource = r.versioner.Resolve(resource)
	}
	if r.cache != nil {
		content, _, err = r.cache.Retrieve(resource)
		if !(err != nil || content == nil || len(content) == 0) {
			return content, false, nil
		}
		cummulative = fmt.Errorf("%s: %w\n", cummulative, err)
		defer func() {
			if repo, _, ver, _ := gop.ProcessRequest(resource); repo != "" && ver != "" && len(content) != 0 {
				_ = r.cache.Cache(resource, content)
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
