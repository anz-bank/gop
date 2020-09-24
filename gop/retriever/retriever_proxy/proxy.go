package retriever_proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joshcarp/gop/gop"
)

type Client struct {
	Proxy  string
	Client *http.Client
}

func New(addr string) Client {
	return Client{Proxy: addr, Client: http.DefaultClient}
}

func (c Client) Retrieve(resource string) ([]byte, bool, error) {
	var resp *http.Response
	var err error
	resp, err = c.Client.Get(c.Proxy + "?resource=" + resource)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.BadRequestError, err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.ProxyReadError, err)
	}
	var obj gop.Object
	if err := json.Unmarshal(bytes, &obj); err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.FileReadError, err)
	}
	return obj.Content, false, nil
}
