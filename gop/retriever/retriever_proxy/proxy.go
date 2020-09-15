package retriever_proxy

import (
	"io/ioutil"
	"net/http"
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
	return bytes, false, nil
}
