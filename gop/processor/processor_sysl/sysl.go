package processor_sysl

import (
	"regexp"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
	"google.golang.org/protobuf/encoding/protojson"
)

var SyslImportRegexStr = `(?:#import.*)|(?:import )(?:\/\/)?(?P<import>.*)`

type Processor struct {
	ImportRegex *regexp.Regexp
}

func (p *Processor) Process(a *pbmod.Object) error {
	if a.Processed != nil && *a.Processed != "" {
		return nil
	}
	withoutImports := p.ImportRegex.ReplaceAllString(a.Content, "")
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
