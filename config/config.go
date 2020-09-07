package config

import (
	"io"
	"io/ioutil"
)

type AppConfig struct {
	Username      string `yaml:"username"`
	Token         string `yaml:"token"`
	CacheLocation string `yaml:"cachelocation"`
}

func ScanIntoString(res *string, file io.Reader) error {
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fin := string(contents)
	*res = fin
	return nil
}
