package retriever_proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/anz-bank/sysl-go/common"
	"github.com/anz-bank/sysl-go/restlib"
	"github.com/anz-bank/sysl-go/validator"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

type Retriever struct {
	AppConfig app.AppConfig
	client    *http.Client
}

func New(appConfig app.AppConfig) Retriever {

	return Retriever{
		AppConfig: appConfig,
		client:    http.DefaultClient,
	}
}

func (a Retriever) Retrieve(repo, resource, version string) (gop.Object, bool, error) {
	req := &gop.GetRequest{
		Resource: path.Join(repo, resource),
		Version:  version,
	}
	ret, err := a.Get(context.Background(), req)
	if err != nil {
		return gop.Object{}, false, err
	}
	return *ret, false, err
}

// Get ...
func (s *Retriever) Get(ctx context.Context, req *gop.GetRequest) (*gop.Object, error) {
	required := []string{}
	var okResponse gop.Object
	u, err := url.Parse(fmt.Sprintf("%s/", s.AppConfig.Proxy))
	if err != nil {
		return nil, app.CreateError(app.BadRequestError, "failed to parse url", err)
	}

	q := u.Query()
	q.Add("resource", req.Resource)

	q.Add("version", req.Version)

	u.RawQuery = fmt.Sprintf("resource=%s&version=%s", req.Resource, req.Version)
	result, err := restlib.DoHTTPRequest(ctx, s.client, "GET", u.String(), nil, required, &okResponse, nil)
	restlib.OnRestResultHTTPResult(ctx, result, err)
	if err != nil {
		return nil, common.CreateError(ctx, common.DownstreamUnavailableError, "call failed: gop <- GET "+u.String(), err)
	}

	if result.HTTPResponse.StatusCode == http.StatusUnauthorized {
		return nil, common.CreateDownstreamError(ctx, common.DownstreamUnauthorizedError, result.HTTPResponse, result.Body, nil)
	}
	OkObjectResponse, ok := result.Response.(*gop.Object)
	if ok {
		valErr := validator.Validate(OkObjectResponse)
		if valErr != nil {
			return nil, common.CreateDownstreamError(ctx, common.DownstreamUnexpectedResponseError, result.HTTPResponse, result.Body, valErr)
		}

		return OkObjectResponse, nil
	}

	return nil, common.CreateDownstreamError(ctx, common.DownstreamUnexpectedResponseError, result.HTTPResponse, result.Body, nil)
}
