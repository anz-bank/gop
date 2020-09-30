package cli

func (r Retriever) Update(old, new string) error {
	return r.versioner.Update(old, new)
}
