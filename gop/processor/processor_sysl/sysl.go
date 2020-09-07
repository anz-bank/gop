package processor_sysl

import (
	"regexp"

	"github.com/joshcarp/pb-mod/app"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
	"google.golang.org/protobuf/encoding/protojson"
)

const default_regex = `(?:#import.*)|(?:import )(?:\/\/)?(?P<import>.*)`

type Processor struct {
	importRegex *regexp.Regexp
}

func New(appConfig app.AppConfig) Processor {
	if appConfig.ImportRegex == "" {
		appConfig.ImportRegex = default_regex
	}
	return Processor{importRegex: regexp.MustCompile(appConfig.ImportRegex)}
}

func (p *Processor) Process(a *pbmod.Object) error {
	if a.Processed != nil && *a.Processed != "" {
		return nil
	}
	withoutImports := p.importRegex.ReplaceAllString(a.Content, "")
	m, err := parse.NewParser().ParseString(withoutImports)
	if err != nil {
		return err
	}
	ma := protojson.MarshalOptions{}
	mb, err := ma.Marshal(m)
	if err != nil {
		return err
	}
	extra := string(mb)
	a.Processed = &extra
	return nil
}
