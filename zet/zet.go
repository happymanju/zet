package zet

import (
	"container/list"
	"slices"
)

type Zet struct {
	Files map[string]*list.List
	Tags  map[string]*list.List
}

func NewZet() *Zet {
	z := Zet{
		Files: make(map[string]*list.List),
		Tags:  make(map[string]*list.List),
	}
	return &z
}

func (z *Zet) UpdateFile(f string, inboundTags []string) {
	u := newUpdateOperation()
	oldTags, ok := z.Files[f]
	if !ok {
		l := list.New()
		for _, tag := range inboundTags {
			l.PushBack(tag)
		}
		z.Files[f] = l
		u.ToAdd = inboundTags
	} else {
		currentOldTag := oldTags.Front()
		for {
			if currentOldTag == nil {
				break
			}
			if !slices.Contains(inboundTags, currentOldTag.Value.(string)) {
				u.ToRemove = append(u.ToRemove, currentOldTag.Value.(string))
			}
			currentOldTag = currentOldTag.Next()
		}

		for _, newTag := range inboundTags {
			tagFound := Find(oldTags, newTag)
			if tagFound != nil {
				u.ToAdd = append(u.ToAdd, newTag)
			}
		}
	}

	// update tag database using the UpdateOperation data
	// make entry for totally new tags, add filename to existing tag entries
	for _, addTag := range u.ToAdd {
		currentTagList, ok := z.Tags[addTag]
		if !ok {
			newTaggedFileList := list.New()
			newTaggedFileList.PushBack(f)
			z.Tags[addTag] = newTaggedFileList
		} else {
			fileFound := Find(currentTagList, f)
			if fileFound == nil {
				currentTagList.PushBack(f)
			}
		}
	}

	// remove filename from the lists in the ToRemove tags
	for _, removeTag := range u.ToRemove {
		currentTagList, ok := z.Tags[removeTag]
		if !ok {
			continue
		} else {
			removeEl := Find(currentTagList, f)
			if removeEl != nil {
				currentTagList.Remove(removeEl)
			}
		}
	}
}
