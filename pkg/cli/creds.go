package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

/* NewTokenMap Returns a tokenstring from a token env var in the form host:token,host:token or from a gitcredentialsvar which points to the gitcredentials location*/
func NewTokenMap(tokensVar, gitCredentialsVar string) (map[string]string, error) {
	tokenmap := make(map[string]string)
	gitCredentials := os.Getenv(gitCredentialsVar)
	if gitCredentials == "" {
		gitCredentials = path.Join(os.Getenv("HOME"), ".git-credentials")
	}
	f, _ := os.Open(gitCredentials)
	if f != nil {
		a, _ := ioutil.ReadAll(f)
		tokenmap, _ = TokensFromGitCredentialsFile(a)
	}
	tokensEnv := os.Getenv(tokensVar)
	var hostTokens []string
	if tokensEnv != "" {
		hostTokens = strings.Split(tokensEnv, ",")
	}
	for _, e := range hostTokens {
		arr := strings.Split(e, ":")
		if len(arr) < 2 {
			return map[string]string{}, fmt.Errorf("TOKENS env var is invalid, should be in form `gita.com:<tokena>,gitb.com:<tokenb>`")
		}
		tokenmap[arr[0]] = arr[1]
	}
	return tokenmap, nil
}

/* TokensFromString returns a map of host:token from a string in form: host:token,host:token */
func TokensFromString(str string) map[string]string {
	hostTokens := strings.Split(str, ",")
	tokenmap := make(map[string]string)
	for _, e := range hostTokens {
		arr := strings.Split(e, ":")
		if len(arr) < 2 {
			return nil
		}
		tokenmap[arr[0]] = arr[1]
	}
	return tokenmap
}

/* TokensFromGitCredentialsFile returns a map of host:token a git credentials file */
func TokensFromGitCredentialsFile(contents []byte) (map[string]string, error) {
	gitCredsRe := regexp.MustCompile(`(?:https:\/\/)(?P<user>.*):(?P<token>.*)(?:@)(?P<host>.*)`)
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	tokenHost := make(map[string]string)
	var token, host string
	for scanner.Scan() {
		for _, match := range gitCredsRe.FindAllStringSubmatch(scanner.Text(), -1) {
			if match == nil {
				continue
			}
			for i, name := range gitCredsRe.SubexpNames() {
				if match[i] != "" {
					switch name {
					case "token":
						token = match[i]
					case "host":
						host = match[i]
					}
				}
			}
			if host != "" && token != "" {
				tokenHost[host] = token
			}
		}
	}
	return tokenHost, nil
}
