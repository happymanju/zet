package zet

import (
	"bufio"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ErrNoMetadata error = errors.New("error: no metadata found")

type FileHashesMap struct {
	FileHashes map[string]string
	isChanged  bool
}

func (fhm *FileHashesMap) Update(p string) error {
	f, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("error in file hash update: could not open file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error in FileHashesMap update: %w", err)
	}

	h, err := GetFNVHash(data)
	if err != nil {
		return fmt.Errorf("error in FileHashesMap update: %w", err)
	}
	//todo
	hexHash := hex.EncodeToString(h)

	v, ok := fhm.FileHashes[p]
	if !ok {
		fhm.FileHashes[p] = hexHash
		fhm.isChanged = true
	} else {
		if v != hexHash {
			fhm.FileHashes[p] = hexHash
		}
	}

	return nil
}

func (fhm *FileHashesMap) Save() error {
	if fhm.isChanged {
		f, err := os.Create("zet.fhm.dat")
		if err != nil {
			return fmt.Errorf("error saving FileHashesMap: %w", err)
		}
		defer f.Close()

		g := gob.NewEncoder(f)
		g.Encode(fhm.FileHashes)
		fhm.isChanged = false
		return nil
	}
	return nil
}

func (fhm *FileHashesMap) Load() error {
	f, err := os.Open("zet.fhm.dat")
	if err != nil {
		return fmt.Errorf("error loading FileHashesMap: %w", err)
	}
	defer f.Close()

	g := gob.NewDecoder(f)
	err = g.Decode(&fhm.FileHashes)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error decoding FileHashesMap file: %w", err)
	}

	return nil
}

func GetFNVHash(data []byte) ([]byte, error) {
	h := fnv.New128()
	_, err := h.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error hashing file: %w", err)
	}

	return h.Sum(nil), nil
}

func BuildIndex() (FileHashesMap, error) {
	fht := FileHashesMap{}

	var walkFunc fs.WalkDirFunc = func(p string, d fs.DirEntry, e error) error {
		if d.IsDir() {
			return filepath.SkipDir
		}
		f, err := os.Open(p)
		if err != nil {
			log.Printf("error opening file %s: %v\n", p, err)
			return nil
		}
		defer f.Close()

		data := []byte{}
		_, err = io.ReadFull(f, data)
		if err != nil {
			log.Printf("error reading from file %s: %v\n", p, err)
			return nil
		}

		h, err := GetFNVHash(data)
		if err != nil {
			log.Printf("error generating hash from file %s data: %v\n", p, err)
			return nil
		}
		stringedFileHash := hex.EncodeToString(h)

		fht.FileHashes[p] = stringedFileHash

		return nil
	}

	err := filepath.WalkDir(".", walkFunc)
	if err != nil {
		return fht, err
	}
	return fht, nil
}

func ParseMetadata(p string) (tags []string, err error) {
	inTags := false

	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", p, err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Scan()
	firstLine := sc.Text()

	if firstLine != "---" {
		return nil, ErrNoMetadata
	}

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "---" {
			break
		}
		var value string

		if !inTags {
			if strings.Contains(line, "tags:") {
				inTags = true
				continue
			} else if strings.Contains(line, ":") {
				inTags = false
				continue
			}
		} else if inTags {
			valueStart := strings.Index(line, "-")
			value = line[valueStart+2:]
			tags = append(tags, value)
			continue
		}
	}

	return tags, nil
}
