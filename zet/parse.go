package zet

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func ParseTags(filename string) ([]string, error) {
	var err error = nil
	foundTags := []string{}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	inMetadata := false
	inTags := false
	for sc.Scan() {
		t := strings.TrimSpace(sc.Text())
		if t == "---" && inMetadata {
			if len(foundTags) == 0 {
				err = errors.New("no tags found in metadata")
			}
			break
		} else if t == "---" && !inMetadata {
			inMetadata = true
		}

		if t == "tags:" {
			inTags = true
			continue
		}

		if inTags {
			newTag := t[strings.Index(t, "-")+2:]
			foundTags = append(foundTags, newTag)
		}
	}

	return foundTags, err
}
