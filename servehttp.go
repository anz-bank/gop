package gop

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	gop3 "github.com/joshcarp/gop/gop"

	"github.com/joshcarp/gop/gop/retriever/retriever_git"

	"github.com/joshcarp/gop/gop/gop_filesystem"
	"github.com/joshcarp/gop/gop/gop_gcs"
	"github.com/spf13/afero"
)

var fs afero.Fs

/* ServeHTTP is a http.HandlerFunc, can be used in deployments like cloud functions*/
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	s, _ := NewGopper(os.Getenv("CacheLocation"), os.Getenv("FsType"))
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
		switch e.Kind {
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

/* GitGopper is an implementation of Gopper that will fall back to a git retriever if the file can't be found in
its data source */
type GitGopper struct {
	gop3.Gopper
	retriever_git.Retriever
}

/* Retrieve attempts to retrieve a file from GitGopper.Gopper, if this fails then it will fall back to using a
git retriever */
func (a GitGopper) Retrieve(resource string) ([]byte, bool, error) {
	if res, cached, err := a.Gopper.Retrieve(resource); err == nil {
		return res, cached, nil
	}
	return a.Retriever.Retrieve(resource)
}

/* NewGopper returns a GitGopper for a config; This Gopper can use an os filesystem, memory filesystem or a gcs bucket*/
func NewGopper(cachelocation, fsType string) (*GitGopper, error) {
	r := GitGopper{}
	switch fsType {
	case "os":
		r.Gopper = gop_filesystem.New(afero.NewOsFs(), "")
	case "mem", "memory", "":
		if fs == nil {
			fs = afero.NewMemMapFs()
		}
		r.Gopper = gop_filesystem.New(fs, "/")
	case "gcs":
		gcs := gop_gcs.New(cachelocation)
		r.Gopper = &gcs
	}
	r.Retriever = retriever_git.New("", "")
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
