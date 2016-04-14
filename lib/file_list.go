package lib

import (
	"errors"
)

// FileListItem - Item in List of files
type FileListItem interface {
	FileName() string
	ID() string
}

// FileList - List of items
type FileList struct {
	list map[string]FileListItem
}

// NewFileList - will create new file list
func NewFileList() *FileList {
	l := &FileList{
		make(map[string]FileListItem),
	}
	return l
}

// Add - will add item into file list
func (l *FileList) Add(i FileListItem) {
	l.list[i.ID()] = i
}

// Find - finds and returns item by id
func (l *FileList) Find(id string) (FileListItem, error) {
	s, ok := l.list[id]
	if !ok {
		return nil, errors.New("Item not found")
	}
	return s, nil
}

func (l *FileList) Delete(id string) {
	delete(l.list, id)
}

// FindByFile - finds if we already have this file prepared
func (l *FileList) FindByFile(filename string) FileListItem {
	for _, s := range l.list {
		if s.FileName() == filename {
			return s
		}
	}
	return nil
}

// GetAll - lists all songs which are in progress
func (l *FileList) GetAll() []FileListItem {
	out := []FileListItem{}
	for _, s := range l.list {
		out = append(out, s)
	}
	return out
}
