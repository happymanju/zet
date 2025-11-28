package main

import (
	"fmt"
	"log"

	"github.com/mpiastucki/zet/zet"
)

func main() {
	z := zet.NewZet()
	t, err := zet.ParseTags("a.md")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("updating")
	z.UpdateFile("a.md", t)
	fmt.Println(z.Files)
	fmt.Println(z.Tags)

}
