package retriever_github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/joshcarp/gop/gop"
)

type Retriever struct {
	token string
}

func New(token string) Retriever {
	return Retriever{token: token}
}

func (a Retriever) Retrieve(resource string) (res []byte, cached bool, err error) {
	var resp *http.Response
	repo, resource, version, err := gop.ProcessRequest(resource)
	if err != nil {
		return nil, false, gop.CreateError(gop.BadRequestError, "Can't process request")
	}

	req, err := url.Parse(fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s",
		strings.ReplaceAll(repo, "github.com/", ""), resource, version))
	cl := http.DefaultClient
	heder := http.Header{}
	heder.Add("accept", "application/vnd.github.v3.raw+json")
	if a.token != "" {
		heder.Add("authorization", "token "+a.token)
	}
	r := &http.Request{
		Method: "GET",
		URL:    req,
		Header: heder,
	}
	resp, err = cl.Do(r)
	if err != nil {
		return res, false, err
	}
	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	return res, false, nil
}
