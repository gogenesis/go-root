// groot-ls recursively dumps the hierarchy tree of a ROOT file
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/binet/go-root/pkg/groot"
)

var fname = flag.String("f", "ntuple.root", "ROOT file to inspect")

func normpath(path []string) string {
	name := strings.Join(path, "/")
	if len(name) > 2 && name[:2] == "//" {
		name = name[1:]
	}
	return name
}

func inspect(dir *groot.Directory, path []string, indent string) {
	name := normpath(path)
	if dir == nil {
		fmt.Printf("err: invalid directory [%s]\n", name)
		return
	}
	keys := dir.Keys()
	nkeys := len(keys)
	str := "|--"
	//fmt.Printf("%s%s -> #%d key(s)\n", indent, name, len(keys))
	for i, k := range keys {
		if i+1 >= nkeys {
			str = "`--"
		}
		fmt.Printf("%s%s %s title='%s' type=%s\n",
			indent, str, k.Name(), k.Title(), k.Class())
		if v, ok := k.Value().(*groot.Directory); ok {
			path := append(path, k.Name())
			inspect(v, path, indent+"    ")
		}
	}
}

func main() {
	fmt.Printf(":: groot-ls ::\n")
	flag.Parse()

	f, err := groot.NewFileReader(*fname)
	if err != nil {
		fmt.Printf("**error**: %v\n", err)
		os.Exit(1)
	}

	if f == nil {
		fmt.Printf("**error**: invalid file pointer\n")
		os.Exit(1)
	}

	fmt.Printf("file: '%s' (version=%v)\n", f.Name(), f.Version())

	dir := f.Dir()
	inspect(dir, []string{"/"}, "")

	fmt.Printf("::bye.\n")
}

// EOF
