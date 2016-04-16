package lib

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

type data struct {
	Name         string
	Songs        []string
	Applications []string
}

// Save - will save data to yaml file
func Save(l LogI, songs *FileList, apps *SoundController, name string) []byte {
	d := &data{
		Name:         name,
		Songs:        songs.FileNames(),
		Applications: apps.AppNames(),
	}
	out, err := yaml.Marshal(d)
	if err != nil {
		l.Error(err.Error())
	}
	f, err := os.Create("last.yml")
	if err != nil {
		l.Error(err.Error())
	} else {
		defer f.Close()
		f.Write(out)
	}

	return out
}

type addableListI interface {
	AddUniq(string, LogI) (bool, error)
}

// Load - will load data from yaml file
func Load(l LogI, songs addableListI, apps addableListI, name string) {
	in, err := ioutil.ReadFile(name)
	if err != nil {
		l.Error(err.Error())
		return
	}
	d := &data{}
	fmt.Println(in)
	yaml.Unmarshal(in, d)
	for _, val := range d.Songs {
		fmt.Println("adding: " + val)
		songs.AddUniq(val, l)
	}
	for _, val := range d.Applications {
		apps.AddUniq(val, l)
	}
}
