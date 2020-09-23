package cli

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"

	"github.com/joshcarp/gop/gop"
)

/* TokensFromString returns a map of host:token from a string in form: host:token,host:token */
func TokensFromString(str string) (map[string]string, error) {
	hostTokens := strings.Split(str, ",")
	tokenmap := make(map[string]string)
	for _, e := range hostTokens {
		arr := strings.Split(e, ":")
		if len(arr) < 2 {
			return nil, gop.UnauthorizedError
		}
		tokenmap[arr[0]] = arr[1]
	}
	return tokenmap, nil
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
				tokenHost[token] = host
			}
		}
	}
	return tokenHost, nil
}
