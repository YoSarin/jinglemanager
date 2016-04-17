package lib

import (
	"errors"
)

// FileListItem - Item in List of files
type FileListItem interface {
	FileName() string
	ID() string
	Remove()
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

// AddUniq - will add unique value
func (l *FileList) AddUniq(filename string, log LogI) (bool, error) {
	f := l.FindByFile(filename)
	if f == nil {
		s, err := NewSong(filename, log)
		if err != nil {
			return false, err
		}
		l.Add(s)
	}
	return true, nil
}

// Find - finds and returns item by id
func (l *FileList) Find(id string) (FileListItem, error) {
	s, ok := l.list[id]
	if !ok {
		return nil, errors.New("Item not found")
	}
	return s, nil
}

// FileNames - will return list of all songs file names
func (l *FileList) FileNames() []string {
	out := make([]string, len(l.list))
	i := 0
	for _, val := range l.list {
		out[i] = val.FileName()
		i++
	}
	return out
}

// Delete - will delete item from list
func (l *FileList) Delete(id string) {
	l.list[id].Remove()
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
func (l *FileList) GetAll() map[string]FileListItem {
	return l.list
}
