package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kazu/fbshelper/fbsparser"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprint(os.Stderr, "usage: fbshelper fbsfile templatefile outputdir")
		return
	}

	fbsfile := os.Args[1]
	tmplate := os.Args[2]
	outDir := os.Args[3]

	bytes, err := ioutil.ReadFile(fbsfile)
	if err != nil {
		fmt.Fprint(os.Stderr, "cannot read file"+fbsfile)
	}

	parser := &fbsparser.Parser{Buffer: string(bytes)}

	parser.Init()
	err = parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error %s err=%s", fbsfile, err)

	}

	parser.Execute()

	for _, info := range parser.Fbs.Structs {

		newSrc, err := info.FromTemplate(tmplate)
		if err == nil {
			output := filepath.Join(outDir, info.Name+".fbshelper.go")
			ioutil.WriteFile(output, []byte(newSrc), 0644)
		} else {
			fmt.Fprint(os.Stderr, err)
		}
	}
}
