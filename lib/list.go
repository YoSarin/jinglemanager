package lib

import (
	"errors"
)

// ListItem - Item in List of files
type ListItem interface {
	ID() string
	Remove()
}

type hasFileName interface {
	FileName() string
}

// UniqueList - List of items
type UniqueList struct {
	list map[string]ListItem
}

// NewUniqueList - will create new file list
func NewUniqueList() *UniqueList {
	l := &UniqueList{
		make(map[string]ListItem),
	}
	return l
}

// Add - will add item into file list
func (l *UniqueList) Add(i ListItem) {
	l.list[i.ID()] = i
}

// AddUniq - will add unique value
func (l *UniqueList) AddUniq(item ListItem, log LogI) (err error) {
	f := l.list[item.ID()]
	if f == nil {
		l.Add(item)
	}
	return
}

// Find - finds and returns item by id
func (l *UniqueList) Find(id string) (ListItem, error) {
	s, ok := l.list[id]
	if !ok {
		return nil, errors.New("Item not found")
	}
	return s, nil
}

// Delete - will delete item from list
func (l *UniqueList) Delete(id string) {
	l.list[id].Remove()
	delete(l.list, id)
}

// FileNames - will return list of all songs file names
func (l *UniqueList) FileNames() []string {
	out := make([]string, len(l.list))
	i := 0
	for _, val := range l.list {
		v, ok := val.(hasFileName)
		if ok {
			out[i] = v.FileName()
			i++
		}
	}
	return out
}

// FindByFile - finds if we already have this file prepared
func (l *UniqueList) FindByFile(filename string) ListItem {
	for _, s := range l.list {
		v, ok := s.(hasFileName)
		if ok {
			if v.FileName() == filename {
				return s
			}
		}
	}
	return nil
}

// GetAll - lists all songs which are in progress
func (l *UniqueList) GetAll() map[string]ListItem {
	return l.list
}
