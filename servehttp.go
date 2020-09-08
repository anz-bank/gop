package gop

import (
	"net/http"
	"os"

	gop3 "github.com/joshcarp/gop/gop"

	"github.com/joshcarp/gop/gop/retriever/retriever_git"

	"github.com/sirupsen/logrus"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gen/pkg/servers/gop"
	gop2 "github.com/joshcarp/gop/gen/pkg/servers/gop"
	"github.com/joshcarp/gop/gop/gop_filesystem"
	"github.com/joshcarp/gop/gop/gop_gcs"
	"github.com/spf13/afero"
)

var fs afero.Fs

/* ServeHTTP is a http.HandlerFunc, can be used in deployments like cloud functions*/
func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	server, err := ServiceHandler(app.AppConfig{
		CacheLocation: os.Getenv("CacheLocation"),
		FsType:        os.Getenv("FsType"),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	server.GetHandler(w, r)
}

/* GitGopper is an implementation of Gopper that will fall back to a git retriever if the file can't be found in
its data source */
type GitGopper struct {
	gop3.Gopper
	retriever_git.Retriever
}

/* Retrieve attempts to retrieve a file from GitGopper.Gopper, if this fails then it will fall back to using a
git retriever */
func (a GitGopper) Retrieve(repo, resource, version string) (gop.Object, bool, error) {
	if res, cached, err := a.Gopper.Retrieve(repo, resource, version); err == nil {
		return res, cached, nil
	}
	return a.Retriever.Retrieve(repo, resource, version)
}

/* NewGopper returns a GitGopper for a config; This Gopper can use an os filesystem, memory filesystem or a gcs bucket*/
func NewGopper(a app.AppConfig) (*GitGopper, error) {
	r := GitGopper{}
	switch a.FsType {
	case "os":
		r.Gopper = gop_filesystem.New(afero.NewOsFs(), a)
	case "mem", "memory", "":
		if fs == nil {
			fs = afero.NewMemMapFs()
		}
		r.Gopper = gop_filesystem.New(fs, a)
	case "gcs":
		gcs := gop_gcs.New(a)
		r.Gopper = &gcs
	}
	r.Retriever = retriever_git.New(a)
	return &r, nil
}

/* ServiceHandler sets up a service handler (Http handler) with a GitGopper */
func ServiceHandler(a app.AppConfig) (*gop2.ServiceHandler, error) {
	g, err := NewGopper(a)
	if err != nil {
		return nil, err
	}
	serve := Server{
		Gopper: g,
	}
	handler, err := gop2.NewServiceHandler(CallBack(), &gop2.ServiceInterface{Get: serve.Get})
	if err != nil {
		return nil, err
	}
	return handler, nil

}
