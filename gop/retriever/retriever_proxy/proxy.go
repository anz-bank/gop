package retriever_proxy

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/joshcarp/gop/gop"
)

type Client struct {
	Proxy string
}

func New(addr string) Client {
	return Client{Proxy: addr}
}

func (c Client) Retrieve(resource string) ([]byte, bool, error) {
	var resp *http.Response
	var err error
	resp, err = http.Get(c.Proxy + "?resource=" + resource)
	if err != nil {
		return nil, false, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	var obj gop.Object
	if err := json.Unmarshal(bytes, &obj); err != nil {
		return nil, false, err
	}
	return obj.Content, false, nil
}
