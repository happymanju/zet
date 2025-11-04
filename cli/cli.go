package cli

import (
	"log"

	"github.com/happymanju/zet/zet"
)

const commandPrompt string = "Enter Command >> "

func Run(args []string) int {
	fhm := zet.FileHashesMap{}
	var err error

	switch args[0] {
	case "rebuild":
		fhm, err = zet.BuildIndex()
		if err != nil {
			log.Printf("error rebuilding index: %v\n", err)
			return 1
		}
	case "search":

	case "update":
	}
	// make file hash table
	// switch: user commands
	// check save flag if saving data structs is needed
	return 0
}
