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
	"strconv"
	"strings"
	"time"
)

func Timestamp() string {
	t := time.Now()
	y, m, d := t.Date()
	currentHour := t.Hour()
	currentMin := t.Minute()

	var stringedMonth string = ""
	var stringedDay string = ""
	var stringedHour string = ""
	var stringedMin string = ""

	if m < 10 {
		stringedMonth = "0" + strconv.Itoa(int(m))
	} else {
		stringedMonth = strconv.Itoa(int(m))
	}

	if d < 10 {
		stringedDay = "0" + strconv.Itoa(d)
	} else {
		stringedDay = strconv.Itoa(int(d))
	}

	if currentHour < 10 {
		stringedHour = fmt.Sprintf("0%d", currentHour)
	} else {
		stringedHour = strconv.Itoa(currentHour)
	}

	if currentMin < 10 {
		stringedMin = fmt.Sprintf("0%d", currentMin)
	} else {
		stringedMin = strconv.Itoa(currentMin)
	}

	return fmt.Sprintf("%v%s%vT%v%v", y, stringedMonth, stringedDay, stringedHour, stringedMin)

}

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

func (z *Zet) NewFile(title string, initialTags string) (titlestring string, err error) {
	tstamp := Timestamp()
	titlestring = strings.ToLower(title)
	titlestring = tstamp + "_" + strings.ReplaceAll(titlestring, " ", "_") + ".md"
	f, err := os.Create(titlestring)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)

	frontmatter := fmt.Sprintf("---\ntags:%s\n---", initialTags)
	_, err = bw.Write([]byte(frontmatter))
	if err != nil {
		return "", err
	}
	err = bw.Flush()
	if err != nil {
		return "", err
	}

	return titlestring, nil
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

func (z *Zet) GetFilesByTitle(keyword string) []string {
	foundFiles := []string{}
	files := maps.Keys(z.FilesToTags)
	for file := range files {
		if strings.Contains(strings.ToLower(file), strings.ToLower(keyword)) {
			foundFiles = append(foundFiles, file)
		}
	}
	return foundFiles
}

func (z *Zet) FindTags(keyword string) []string {
	foundTags := []string{}
	for _, tag := range z.Tags {
		if strings.Contains(tag, keyword) {
			foundTags = append(foundTags, tag)
		}
	}
	return foundTags
}
