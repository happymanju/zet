package zet

import "container/list"

func Find(l *list.List, s string) *list.Element {
	var foundEl *list.Element = nil
	currentEl := l.Front()
	for {
		if currentEl == nil {
			break
		}
		if currentEl.Value.(string) == s {
			foundEl = currentEl
			break
		}
		currentEl = currentEl.Next()
	}

	return foundEl
}
