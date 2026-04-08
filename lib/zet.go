package zet

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func ParseTags(fp string) (tags []string, err error) {
	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	delim := "---"

	sc := bufio.NewScanner(f)

	sc.Scan()
	if sc.Text() != delim {
		return nil, errors.New("malformed frontmatter")
	}

	for sc.Scan() {
		t := sc.Text()
		if t == delim {
			break
		}
		if strings.Contains(t, "tags:") {
			idx := strings.Index(t, ":")
			tags = strings.Split(t[idx+1:], ",")
		}
	}
	if tags == nil {
		return nil, fmt.Errorf("No tags found in: %s\n", fp)
	}
	return tags, nil
}

// return a list of files by tag
// return a list of all current tags, sorted
// search tags by keyword
// rebuild cache
// add a file to cache
// detect a missing file and remove it from cache
type Zet struct {
	FilesToTags map[string][]string
	Tags        []string
}

func NewZet() *Zet {
	z := Zet{
		FilesToTags: make(map[string][]string),
		Tags:        []string{},
	}
	return &z
}

func (z *Zet) Load(fp string) error {
	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err = dec.Decode(z)
	if err != nil {
		return err
	}
	return nil
}

func (z *Zet) Save(fp string) error {
	f, err := os.Create(filepath.Clean(fp))
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(z)
	if err != nil {
		return err
	}
	return nil
}

func (z *Zet) AddFile(fp string) error {
	_, ok := z.FilesToTags[fp]
	if !ok {
		newtags, err := ParseTags(fp)
		if err != nil {
			return err
		}
		z.FilesToTags[fp] = newtags
		for _, newTag := range newtags {
			if !slices.Contains(z.Tags, newTag) {
				z.Tags = append(z.Tags, newTag)
			}
		}
		slices.Sort(z.Tags)
	}

	return nil
}

func (z *Zet) RemoveFile(fp string) {
	_, ok := z.FilesToTags[fp]
	if ok {
		delete(z.FilesToTags, fp)
	} else {
		return
	}
	z.PruneTags()

}

func (z *Zet) UpdateFile(fp string) error {
	var newtags []string
	newtags, err := ParseTags(fp)
	if err != nil {
		return err
	}
	for _, newTag := range newtags {
		if !slices.Contains(z.Tags, newTag) {
			z.Tags = append(z.Tags, newTag)
		}
	}
	z.FilesToTags[fp] = newtags
	slices.Sort(z.Tags)
	z.PruneTags()

	return nil
}

func (z *Zet) PruneTags() {
	prunedTags := []string{}
	for _, tag := range z.Tags {
		found := false
		files := maps.Keys(z.FilesToTags)
		for file := range files {
			found = slices.Contains(z.FilesToTags[file], tag)
		}
		if found {
			prunedTags = append(prunedTags, tag)
		}
	}
	z.Tags = prunedTags
}

func (z *Zet) GetFilesByTag(tag string) []string {
	foundFiles := []string{}

	keys := maps.Keys(z.FilesToTags)
	for file := range keys {
		if slices.Contains(z.FilesToTags[file], tag) {
			foundFiles = append(foundFiles, file)
		}
	}
	return foundFiles
}
