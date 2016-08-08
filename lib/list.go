package lib

import (
	"errors"
	"sync"
)

// ListItem - Item in List of files
type ListItem interface {
	ID() string
	OnRemove()
}

type hasFileName interface {
	FileName() string
}

// UniqueList - List of items
type UniqueList struct {
	list map[string]ListItem
	m    *sync.Mutex
}

// NewUniqueList - will create new file list
func NewUniqueList() *UniqueList {
	l := &UniqueList{
		make(map[string]ListItem),
		&sync.Mutex{},
	}
	return l
}

// Add - will add item into file list
func (l *UniqueList) Add(i ListItem) {
	l.m.Lock()
	l.list[i.ID()] = i
	l.m.Unlock()
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
	l.m.Lock()
	s, ok := l.list[id]
	l.m.Unlock()
	if !ok {
		return nil, errors.New("Item not found")
	}
	return s, nil
}

// Delete - will delete item from list
func (l *UniqueList) Delete(id string) {
	l.m.Lock()
	l.list[id].OnRemove()
	delete(l.list, id)
	l.m.Unlock()
}

// FileNames - will return list of all songs file names
func (l *UniqueList) FileNames() []string {
	out := make([]string, len(l.list))
	i := 0
	l.m.Lock()
	for _, val := range l.list {
		v, ok := val.(hasFileName)
		if ok {
			out[i] = v.FileName()
			i++
		}
	}
	l.m.Unlock()
	return out
}

// FindByFile - finds if we already have this file prepared
func (l *UniqueList) FindByFile(filename string) ListItem {
	l.m.Lock()
	for _, s := range l.list {
		v, ok := s.(hasFileName)
		if ok {
			if v.FileName() == filename {
				l.m.Unlock()
				return s
			}
		}
	}
	l.m.Unlock()
	return nil
}

// GetAll - lists all songs which are in progress
func (l *UniqueList) GetAll() map[string]ListItem {
	return l.list
}
