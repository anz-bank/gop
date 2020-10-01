package retrievertests

import "github.com/joshcarp/gop/gop"

var Tests = map[string]string{
	"github.com/joshcarp/gop/.gitignore@104c2af996a36b532113194dbf359ed128ad0b32":     "/default.db\n/out",
	"github.com/joshcarp/gop/.gitignore@main":                                         "/default.db\n/out",
	"github.com/joshcarp/sysl-1/sysl.sysl@main":                                       "import sysl-1.sysl\n\n_:\n    ...\n",
	"github.com/joshcarp/sysl-1/sysl-1.sysl@main":                                     "import //github.com/joshcarp/sysl-2/sysl-2.sysl\n\nApp:\n    @package = \"App\"\n    /address:\n        GET(ok <: foo):\n            return ok <: sequence of App.foo\n    !type foo:\n        this <: string\n        that <: int\nBar:\n    /bar/{id <: string}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /sameappref/{id <: tar}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /ref/{id <: App.foo}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /address:\n        GET?street=string:\n            return ok <: sequence of foo\n    !type foo:\n        this <: string\n        that <: Bar.tar\n    !type tar:\n        this <: string\n        that <: int\n",
	"github.com/joshcarp/sysl-1/sysl.sysl@ee5c8cd2b97ba24226cb86556c17ad4f852915f2":   "import sysl-1.sysl\n\n_:\n    ...\n",
	"github.com/joshcarp/sysl-1/sysl-1.sysl@911c664b22f5b8dedb7f1f0554ae3ea77085eaac": "import //github.com/joshcarp/sysl-2/sysl-2.sysl@main\n\nApp:\n    @package = \"App\"\n    /address:\n        GET(ok <: foo):\n            return ok <: sequence of App.foo\n    !type foo:\n        this <: string\n        that <: int\nBar:\n    /bar/{id <: string}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /sameappref/{id <: tar}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /ref/{id <: App.foo}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /address:\n        GET?street=string:\n            return ok <: sequence of foo\n    !type foo:\n        this <: string\n        that <: Bar.tar\n    !type tar:\n        this <: string\n        that <: int\n",
	"github.com/joshcarp/sysl-1/sysl-1.sysl@v1.0.0":                                   "import //github.com/joshcarp/sysl-2/sysl-2.sysl@main\n\nApp:\n    @package = \"App\"\n    /address:\n        GET(ok <: foo):\n            return ok <: sequence of App.foo\n    !type foo:\n        this <: string\n        that <: int\nBar:\n    /bar/{id <: string}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /sameappref/{id <: tar}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /ref/{id <: App.foo}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /address:\n        GET?street=string:\n            return ok <: sequence of foo\n    !type foo:\n        this <: string\n        that <: Bar.tar\n    !type tar:\n        this <: string\n        that <: int\n",
}

var GithubRequestPaths = map[string]string{
	/* These mock the "repo/content" endpoint */
	"/repos/joshcarp/gop/contents/.gitignore?ref=104c2af996a36b532113194dbf359ed128ad0b32":     "/default.db\n/out",
	"/repos/joshcarp/gop/contents/.gitignore?ref=main":                                         "/default.db\n/out",
	"/repos/joshcarp/sysl-1/contents/sysl.sysl?ref=main":                                       "import sysl-1.sysl\n\n_:\n    ...\n",
	"/repos/joshcarp/sysl-1/contents/sysl-1.sysl?ref=main":                                     "import //github.com/joshcarp/sysl-2/sysl-2.sysl\n\nApp:\n    @package = \"App\"\n    /address:\n        GET(ok <: foo):\n            return ok <: sequence of App.foo\n    !type foo:\n        this <: string\n        that <: int\nBar:\n    /bar/{id <: string}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /sameappref/{id <: tar}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /ref/{id <: App.foo}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /address:\n        GET?street=string:\n            return ok <: sequence of foo\n    !type foo:\n        this <: string\n        that <: Bar.tar\n    !type tar:\n        this <: string\n        that <: int\n",
	"/repos/joshcarp/sysl-1/contents/sysl.sysl?ref=ee5c8cd2b97ba24226cb86556c17ad4f852915f2":   "import sysl-1.sysl\n\n_:\n    ...\n",
	"/repos/joshcarp/sysl-1/contents/sysl-1.sysl?ref=911c664b22f5b8dedb7f1f0554ae3ea77085eaac": "import //github.com/joshcarp/sysl-2/sysl-2.sysl@main\n\nApp:\n    @package = \"App\"\n    /address:\n        GET(ok <: foo):\n            return ok <: sequence of App.foo\n    !type foo:\n        this <: string\n        that <: int\nBar:\n    /bar/{id <: string}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /sameappref/{id <: tar}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /ref/{id <: App.foo}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /address:\n        GET?street=string:\n            return ok <: sequence of foo\n    !type foo:\n        this <: string\n        that <: Bar.tar\n    !type tar:\n        this <: string\n        that <: int\n",
	"/repos/joshcarp/sysl-1/contents/sysl.sysl?ref=v1.0.0":                                     "import sysl-1.sysl\n\n_:\n    ...\n",
	"/repos/joshcarp/sysl-1/contents/sysl-1.sysl?ref=v1.0.0":                                   "import //github.com/joshcarp/sysl-2/sysl-2.sysl@main\n\nApp:\n    @package = \"App\"\n    /address:\n        GET(ok <: foo):\n            return ok <: sequence of App.foo\n    !type foo:\n        this <: string\n        that <: int\nBar:\n    /bar/{id <: string}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /sameappref/{id <: tar}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /ref/{id <: App.foo}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /address:\n        GET?street=string:\n            return ok <: sequence of foo\n    !type foo:\n        this <: string\n        that <: Bar.tar\n    !type tar:\n        this <: string\n        that <: int\n",

	/* These mock the "commits" endpoint */
	"/repos/joshcarp/gop/commits/test":  "dad0c54cae43ea40f3f1b5063af680ed4521eab2",
	"/repos/joshcarp/gop/commits/test2": "dad0c54cae43ea40f3f1b5063af680ed4521eab2",

	"/repos/joshcarp/sysl-1/contents/sysl-1.sysl?ref=e52c640e41ba2cc918e4f2dda3a2cfeb4768b075":           "import //github.com/joshcarp/sysl-2/sysl-2.sysl\n\nApp:\n    @package = \"App\"\n    /address:\n        GET(ok <: foo):\n            return ok <: sequence of App.foo\n    !type foo:\n        this <: string\n        that <: int\nBar:\n    /bar/{id <: string}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /sameappref/{id <: tar}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /ref/{id <: App.foo}:\n        GET:\n            App <- GET /address\n            return ok <: foo\n    /address:\n        GET?street=string:\n            return ok <: sequence of foo\n    !type foo:\n        this <: string\n        that <: Bar.tar\n    !type tar:\n        this <: string\n        that <: int\n",
	"/repos/joshcarp/sysl-1/contents/sysl_modules/sysl.mod?ref=e52c640e41ba2cc918e4f2dda3a2cfeb4768b075": "\ndirect:\n- pattern: github.com/joshcarp/sysl-2\n  resolve: github.com/joshcarp/sysl-2@0e19b891da4ea38e82910a5ef3dc24692ab711ce",
	"/repos/joshcarp/sysl-1/commits/main":                                                      "e52c640e41ba2cc918e4f2dda3a2cfeb4768b075",
	"/repos/joshcarp/sysl-2/commits/0e19b891da4ea38e82910a5ef3dc24692ab711ce":                  "0e19b891da4ea38e82910a5ef3dc24692ab711ce",
	"/repos/joshcarp/sysl-2/contents/sysl-2.sysl?ref=0e19b891da4ea38e82910a5ef3dc24692ab711ce": "\nsdifnscfsdhifs:\n    /address:\n        GET(ok <: foo):\n            return ok <: sequence of App.foo\n    !type foo:\n        this <: string\n        that <: int",
}

type retriever struct {
	content map[string]string
}

func (r retriever) Retrieve(resource string) (content []byte, cached bool, err error) {
	if val, ok := r.content[resource]; ok {
		return []byte(val), cached, nil
	}
	return nil, false, gop.FileNotFoundError
}

func (r retriever) Cache(resource string, content []byte) (err error) {
	r.content[resource] = string(content)
	return nil
}

func New(a map[string]string) retriever {
	return retriever{content: a}
}
