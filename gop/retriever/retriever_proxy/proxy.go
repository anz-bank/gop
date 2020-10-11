package retriever_proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/joshcarp/gop/gop"
)

type Client struct {
	Proxy  string
	Client *http.Client
	Header http.Header
}

func New(addr string) Client {
	return Client{Proxy: addr, Client: &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}}
}

func (c *Client) SetHeader(header http.Header) {
	c.Header = header
}

func (c Client) Retrieve(resource string) ([]byte, bool, error) {
	var resp *http.Response
	var err error
	req := &http.Request{Header: c.Header}
	rawurl, err := url.Parse(c.Proxy + "?resource=" + resource)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.BadRequestError, err)
	}
	req.URL = rawurl
	resp, err = c.Client.Do(req)
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
