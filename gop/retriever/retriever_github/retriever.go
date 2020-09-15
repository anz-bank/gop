package retriever_github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gop"
)

type Retriever struct {
	AppConfig app.AppConfig
}

func New(appConfig app.AppConfig) Retriever {
	return Retriever{AppConfig: appConfig}
}

func (a Retriever) Retrieve(repo, resource, version string) (res gop.Object, cached bool, err error) {
	var resp *http.Response
	req := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s",
		strings.ReplaceAll(repo, "github.com/", ""), version, resource)
	resp, err = http.Get(req)
	if err != nil {
		return res, false, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, false, err
	}
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return res, false, err
	}
	return res, false, nil
}
