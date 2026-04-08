package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	zet "github.com/happymanju/zet/lib"
)

const defaultCache string = "zet_cache.bin"

func Run(args []string) int {
	z := zet.NewZet()
	z.Load(defaultCache)
	cacheModified := false

	switch args[0] {
	case "add":
		err := z.AddFile(args[1])
		if err != nil {
			log.Println(err)
			return 1
		}
		cacheModified = true
	case "update":
		err := z.UpdateFile(args[1])
		if err != nil {
			log.Println(err)
			return 1
		}
		cacheModified = true
	case "del":
		z.RemoveFile(args[1])
		cacheModified = true
	case "-h":
		fmt.Println("add, update, del, tag -a -s -f")

	case "tag":
		sc := bufio.NewScanner(os.Stdin)
		switch args[1] {
		case "-a":
			for _, v := range z.Tags {
				fmt.Println(v)
			}
		case "-s":
			fmt.Print("keyword search term >> ")
			sc.Scan()
			t := sc.Text()
			for _, tag := range z.Tags {
				if strings.Contains(tag, t) {
					fmt.Println(tag)
				}
			}

		case "-f":
			fmt.Print("retrieve files by tag >> ")
			sc.Scan()
			t := sc.Text()
			foundFiles := z.GetFilesByTag(t)
			if foundFiles != nil {
				for _, file := range foundFiles {
					fmt.Println(file)
				}
			}
		}
	}

	if cacheModified {
		err := z.Save(defaultCache)
		if err != nil {
			log.Println(err)
			return 1
		}
	}
	return 0
}
