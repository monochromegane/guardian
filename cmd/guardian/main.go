package main

import (
	"os"

	"github.com/monochromegane/guardian"
)

func main() {
	os.Exit(guardian.Run(os.Args[1:]))
}
