package server

import (
	"path"
	"strings"
)

func processRequest(resource string) (string, string) {
	parts := strings.Split(resource, "/")
	if len(parts) < 3 {
		return "", ""
	}
	repo := path.Join(parts[0], parts[1], parts[2])
	relresource := path.Join(parts[3:]...)
	return repo, relresource
}

//func TestFindImport(t *testing.T) {
//	tests := map[string][]string{
//		`
//#import notimported
//import a.sysl
//import b.sysl`: {"a.sysl", "b.sysl"},
//	}
//	for in, out := range tests {
//		t.Run(in, func(t *testing.T) {
//			a := findImports(syslimportRegex, []byte(in))
//			require.Equal(t, out, a)
//		})
//	}
//}

///* re is the regex that matches the import statements, and  */
//func findImports(importRegex string, file []byte) []string {
//	var re = regexp.MustCompile(importRegex)
//	scanner := bufio.NewScanner(bytes.NewReader(file))
//	var imports []string
//	for scanner.Scan() {
//		for _, match := range re.FindAllStringSubmatch(scanner.Text(), -1) {
//			if match == nil {
//				continue
//			}
//			for i, name := range re.SubexpNames() {
//				if name == "import" && match[i] != "" {
//					if path.Ext(match[i]) != ".sysl" {
//						match[i] += ".sysl"
//					}
//					imports = append(imports, match[i])
//				}
//			}
//		}
//	}
//	return imports
//}
