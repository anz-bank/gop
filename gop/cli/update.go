package cli

import "github.com/joshcarp/gop/gop"

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
	}
	return r.versioner.Update(repo)
}

func (r Retriever) Update(repo string) error {
	return r.versioner.Update(repo)
}

func (r Retriever) UpdateAll() error {
	return r.versioner.UpdateAll()
}

func (r Retriever) Get(repo string) error {
	a, _, c := gop.ProcessRepo(repo)
	return r.versioner.UpdateTo(a, gop.CreateResource(a, "", c))
}
