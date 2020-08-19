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
