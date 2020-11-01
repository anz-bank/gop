package cli

import "github.com/anz-bank/gop/pkg/gop"

/* Command is meant to be used in a cli tool with cmd and repo args */
func (r Retriever) Command(cmd, repo string) error {
	switch cmd {
	case "get":
		if _, resource, _ := gop.ProcessRepo(repo); resource != "" {
			_, _, err := r.Retrieve(repo)
			return err
		} else {
			return r.Get(repo)
		}
	case "update":
		if repo == "" {
			return r.UpdateAll()
		} else {
			return r.Update(repo)
		}
	case "init":
		return r.Init()
	}
	return r.Update(repo)
}
