package retriever_proxy

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/joshcarp/gop/gop"
)

type Client struct {
	Proxy string
}

func New(addr string) Client {
	return Client{Proxy: addr}
}

func (c Client) Retrieve(repo, resource, version string) (res gop.Object, cached bool, err error) {
	var resp *http.Response
	resp, err = http.Get(c.Proxy + "?resource=" + path.Join(repo, resource) + "@" + version)
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
