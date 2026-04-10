package cli

import (
	"bufio"
	"fmt"
	"io/fs"
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

	case "new":
		sc := bufio.NewScanner(os.Stdin)

		fmt.Print("new title >> ")
		sc.Scan()
		title := sc.Text()

		fmt.Println("")
		fmt.Print("initial tags >> ")
		sc.Scan()
		initialTags := sc.Text()

		newFile, err := z.NewFile(title, initialTags)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(newFile)
		err = z.AddFile(newFile)
		if err != nil {
			log.Println(err)
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

	case "rebuild":
		dirEntries, err := fs.ReadDir(os.DirFS("."), ".")
		if err != nil {
			log.Println(err)
			break
		}
		z.Tags = []string{}
		z.FilesToTags = map[string][]string{}
		for _, entry := range dirEntries {
			if !entry.IsDir() && entry.Name() != defaultCache {
				z.AddFile(entry.Name())
			}
		}

	case "search":
		foundTags := z.FindTags(args[1])
		foundFiles := z.GetFilesByTitle(args[1])
		foundFilesByTag := z.GetFilesByTag(args[1])

		fmt.Println("Found Tags:")
		for _, v := range foundTags {
			fmt.Println(v)
		}

		fmt.Println("Found File Titles:")
		for _, file := range foundFiles {
			fmt.Println(file)
		}

		fmt.Println("Found Files by Tag:")
		for _, file := range foundFilesByTag {
			fmt.Println(file)
		}

	case "-h":
		fmt.Println("add, update, del, search, rebuild, tag -a -s -f")

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
			if foundFiles == nil {
				break
			}
			for _, file := range foundFiles {
				fmt.Println(file)
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
