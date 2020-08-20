# sysl-mod
modules for sysl

## Requirements
- Ability to easily import from and git service (github/gitlab/bitbucket)
    - Must be platform agnostic
- Ability to import from multiple versions of a specification without name clashes
- Ability to be language agnostic
    - Can't assume that the repository has any other code in it
    - Current solution requires a go.mod file to exist in target repo

## Nice to haves
- Fast
    - Sysl parser is slow at the moment, large projects take minutes to parse

## Inspirations
- go modules: https://github.com/gomods/athens

## Ideas
- a grpc service that stores minimal versions of repos (just .sysl files and their protobuf equivilents)
Because sysl module objects are already protobuf objects, returning an already compiled sysl module makes sense, and will be extremely fast

- dependencies can be found from the `import` statements at the start of a sysl file

- .sysl files would be stored along side .sysl.pb with the raw protobuf bytes 

- .sysl.pb files would only have the compiled protobuf from the .sysl file; this would mitigate the risk of reimporting multiple of the same modules if they're imported in a loop

- A grpc gateway could be provided for a json object of any endpoint to ensure compatibility

## Example

given a sysl file:
```
// example.sysl
import //github.com/blah/bah@ver
// application bar is imported

foo:
    ep:
        bar@ver <- ep // bar is mapped to bar@ver to avoid clashes

   
```

The sequence diagram for this compilation would be: 

<img src="http://www.plantuml.com/plantuml/png/dP31QiCm44Jl-eezbpY-cr82EVGKw2_8rX89LbgpMfM6qd-lxIBHqb93RdkOcRSpfwnMj4GoPmgO5BedkB1x4Nwx3N15XRw_1lLbF4uS-v6ixqVhZR6a4DaLGfZivD4P06ZMDUOhS011BPBW8Tyo7I-RPTCsO5EUro1u_wuym2jA3fmE4EBCeld386MixCJwDt-9VTx-_g_53yjZrpMBuCo_2hLxWAmikAjQM7EW-gl1vEhAAwiAwqBxnUzUxBIWIwHF">



Within the sysl proxy filesystem/database the following would exist:
```
proxy.sysl.io/
              github.com/
                         blah/bah/
                                 @ver/
                                      example.sysl // raw sysl file
                                      example.sysl.pb // compiled protobuf bytes (only of the sysl file, not of the imports)
                                      example.sysl.imports.json // a json list of modules imported by this file
```
Now next time this module is requested, the proxy can just return example.sysl.pb (a *sysl.Module) and the *sysl.Modules of `example.sysl.imports.json`

## What about arrai/yaml/json?

Same thing, instead of returning a *sysl.Module, we can just return any other protobuf message, including a raw string of arrai bytes
