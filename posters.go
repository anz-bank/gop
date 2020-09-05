package main

import (
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
	"google.golang.org/protobuf/encoding/protojson"
)

type post func(pre *pbmod.KeyValue) (err error)

func processSysl(a *pbmod.KeyValue) error {
	if *a.Extra != "" {
		return nil
	}
	m, err := parse.NewParser().ParseString(a.Value)
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
