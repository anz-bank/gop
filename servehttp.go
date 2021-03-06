package gop

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/anz-bank/gop/pkg/cli"
	"github.com/anz-bank/gop/pkg/retrievers/github"

	gop3 "github.com/anz-bank/gop/pkg/gop"

	"github.com/anz-bank/gop/pkg/goppers/filesystem"
	"github.com/anz-bank/gop/pkg/goppers/gcs"
	"github.com/spf13/afero"
)

var fs afero.Fs

/* ServeHTTP is a http.HandlerFunc, can be used in deployments like cloud functions*/
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	s, _ := NewGopper(os.Getenv("CacheLocation"), os.Getenv("FsType"), "GIT_TOKEN")
	defer func() {
		HandleErr(w, err)
	}()
	reqestedResource := r.URL.Query().Get("resource")

	var res []byte
	var cached bool
	res, cached, err = s.Retrieve(reqestedResource)
	if err != nil || res == nil {
		return
	}
	if !cached {
		if err := s.Cache(reqestedResource, res); err != nil {
			return
		}
	}
	b, err := json.Marshal(gop3.Object{Content: res, Resource: reqestedResource})
	if err != nil {
		return
	}
	if _, err := w.Write(b); err != nil {
		log.Println(err)
	}
}

func HandleErr(w http.ResponseWriter, err error) {
	var httpCode int
	var desc string
	if err == nil {
		return
	}
	log.Println(err)
	switch e := err.(type) {
	case gop3.Error:
		desc = e.String()
		switch e {
		case gop3.BadRequestError:
			httpCode = 400
		case gop3.UnauthorizedError:
			httpCode = 401
		case gop3.TimeoutError:
			httpCode = 408
		case gop3.CacheAccessError, gop3.CacheWriteError:
			httpCode = 503
		case gop3.CacheReadError, gop3.FileNotFoundError:
			httpCode = 404
		default:
			httpCode = 500
		}
	default:
		httpCode = 500
		desc = "Unknown"
	}
	http.Error(w, desc, httpCode)
}

/* GopperService is an implementation of Gopper that will fall back to a git retriever if the file can't be found in
its data source */
type GopperService struct {
	gop3.Gopper
	gop3.Retriever
	gop3.Resolver
}

/* Retrieve attempts to retrieve a file from GopperService.Gopper, if this fails then it will fall back to using a
git retriever */
func (a GopperService) Retrieve(resource string) ([]byte, bool, error) {
	if res, cached, err := a.Gopper.Retrieve(resource); err == nil {
		return res, cached, nil
	}
	res, cached, err := a.Retriever.Retrieve(resource)
	if err != nil {
		return nil, false, err
	}
	if !cached {
		_ = a.Gopper.Cache(resource, res)
	}
	return res, cached, err
}

/* NewGopper returns a GopperService for a config; This Gopper can use an os filesystem, memory filesystem or a gcs bucket*/
func NewGopper(cachelocation, fsType, tokenEnv string) (*GopperService, error) {
	r := GopperService{}
	switch fsType {
	case "os":
		r.Gopper = filesystem.New(afero.NewOsFs(), "")
	case "mem", "memory", "":
		if fs == nil {
			fs = afero.NewMemMapFs()
		}
		r.Gopper = filesystem.New(fs, "/")
	case "gcs":
		gcs := gcs.New(cachelocation)
		r.Gopper = &gcs
	}
	token, _ := cli.NewTokenMap(tokenEnv, "")
	r.Retriever = github.New(token)
	return &r, nil
}

func CORSEnabledFunction(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
