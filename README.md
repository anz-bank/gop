# GOP
Git Object Proxy
GOP, is a generic tool that allows for git objects to be retrieved, processed and cached. It can be used as a generic package or can be started as a service.
See [design/revision2.md](design/revision2.md) for a detailed explanation about how GOP works and is implemented

### Requirements
    - go version 1.14

### Start the server
`go run . config.yaml`
- This will start a server on localhost:8080
#### Invoke the service
`curl localhost:8080/?resource=<resource>&version=<hash>`
Where:
 - `<resource>` is the repo + path of the resource (eg `github.com/joshcarp/gop/pbmod.sysl`)
 - `<hash>` is the hash of the commit to retrieve
 