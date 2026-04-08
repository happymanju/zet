package main

import (
	"os"

	"github.com/happymanju/zet/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
