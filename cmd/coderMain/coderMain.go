package main

import (
	"fmt"
	"l4/pkg/coder"
	"os"

	"github.com/ftrvxmtrx/tga"
)

func main() {
	fmt.Println(os.Getwd())
	var path string
	if len(os.Args) < 2 {
		path = "data/input/testy4/example1.tga"
	} else {
		path = os.Args[1]
	}

	file, err := os.Open(path)

	if err != nil {
		os.Exit(1)
	}

	img, err2 := tga.Decode(file)

	if err2 != nil {
		os.Exit(1)
	}

	c := coder.Coder_createCoder(img)
	c.Coder_run()

}
