package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TimothyYe/glance/lib"

	"github.com/TimothyYe/glance/core"
	"github.com/TimothyYe/glance/reader"
)

var (
	r       reader.Reader
	Version string
)

func main() {
	lines := flag.Int("n", 1, "number of lines to display at once")
	showVersion := flag.Bool("v", false, "display version")
	flag.Parse()

	if *showVersion {
		lib.DisplayVersion(Version)
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Println("Please input the filename")
		os.Exit(1)
	}

	if *lines < 1 {
		fmt.Println("line count must be greater than 0")
		os.Exit(1)
	}

	filePath := flag.Arg(0)
	absPath, err := filepath.Abs(filePath)
	if err == nil {
		filePath = absPath
	}
	ext := strings.ToUpper(filepath.Ext(filePath))

	// create reader
	switch ext {
	case ".TXT":
		r = reader.Reader(reader.NewTxtReader())
	case ".EPUB":
		r = reader.Reader(reader.NewEpubReader())
	default:
		fmt.Println("Unsupported file format!")
		os.Exit(1)
	}

	if err := r.Load(filePath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	progress, err := lib.LoadProgress()
	if err == nil {
		if pos, ok := progress.Books[filePath]; ok {
			r.Goto(pos)
		}
	}

	core.Init(r, *lines, func(pos int) {
		store, err := lib.LoadProgress()
		if err != nil {
			return
		}

		store.Books[filePath] = pos
		_ = lib.SaveProgress(store)
	}, filePath)
}
