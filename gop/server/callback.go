package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/anz-bank/sysl-go/common"

	"github.com/joshcarp/gop/app"

	"github.com/anz-bank/sysl-go/validator"
	"github.com/go-chi/chi"
)

func CallBack() Callback {
	return Callback{
		UpstreamTimeout:   120 * time.Second,
		DownstreamTimeout: 120 * time.Second,
		RouterBasePath:    "/",
		UpstreamConfig:    Config{},
	}
}

type Callback struct {
	UpstreamTimeout   time.Duration
	DownstreamTimeout time.Duration
	RouterBasePath    string
	UpstreamConfig    validator.Validator
}

type Config struct{}

func (c Config) Validate() error {
	return nil
}

func (g Callback) AddMiddleware(ctx context.Context, r chi.Router) {
}

func (g Callback) BasePath() string {
	return g.RouterBasePath
}

func (g Callback) Config() interface{} {
	return g.UpstreamConfig
}

func (g Callback) HandleError(ctx context.Context, w http.ResponseWriter, kind app.Kind, message string, cause error) {
	se := app.CreateError(kind, message, cause)
	g.MapError(ctx, se).WriteError(ctx, w)
}

func (g Callback) DownstreamTimeoutContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, g.DownstreamTimeout)
}

func (g Callback) MapError(ctx context.Context, err error) *common.HTTPError {
	httpErr := MapError(ctx, err)
	return &httpErr
}

func MapError(ctx context.Context, err error) common.HTTPError {
	var (
		httpCode int
		desc     string
	)

	switch e := err.(type) {
	case app.Error:
		desc = e.String()
		switch e.Kind {
		case app.BadRequestError:
			httpCode = 400
		case app.UnauthorizedError:
			httpCode = 401
		case app.TimeoutError:
			httpCode = 408
		case app.CacheAccessError, app.CacheWriteError:
			httpCode = 503
		case app.CacheReadError, app.FileNotFoundError:
			httpCode = 404
		default:
			httpCode = 500
		}
	case *common.ServerError:
		switch e := e.Unwrap().(type) {
		case app.Error:
			desc = e.String()
			switch e.Kind {
			case app.BadRequestError:
				httpCode = 400
			case app.UnauthorizedError:
				httpCode = 401
			case app.TimeoutError:
				httpCode = 408
			case app.CacheAccessError, app.CacheWriteError:
				httpCode = 503
			case app.CacheReadError, app.FileNotFoundError:
				httpCode = 404
			default:
				httpCode = 500
			}
		default:
			httpCode = 500
			desc = "Unknown"
		}
	default:
		httpCode = 500
		desc = "Unknown"
	}
	return common.HTTPError{
		HTTPCode:    httpCode,
		Code:        strconv.Itoa(httpCode),
		Description: desc,
	}
}
