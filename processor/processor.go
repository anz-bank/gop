package processor

import (
	"regexp"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
	"google.golang.org/protobuf/encoding/protojson"
)

type Processor interface {
	Processor(pre *pbmod.Module) (err error)
}

var SyslImportRegexStr = `(?:#import.*)|(?:import )(?:\/\/)?(?P<import>.*)`

type ProcessorSysl struct {
	SyslimportRegex *regexp.Regexp
}

func (p *ProcessorSysl) Processor(a *pbmod.Module) error {
	if a.Extra != nil && *a.Extra != "" {
		return nil
	}
	withoutImports := p.SyslimportRegex.ReplaceAllString(a.Value, "")
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
	a.Extra = &extra
	return nil
}
