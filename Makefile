  
SYSLFILE = gop.sysl
APPS = gop

-include local.mk
include codegen.mk

run:
	go run . config.yaml