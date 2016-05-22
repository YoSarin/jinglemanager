package lib

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type data struct {
	Tournament   *Tournament
	Applications []string
	Jingles      []*JingleStorage
}

// Save - will save data to yaml file
func (c *Context) Save() []byte {
	if c.Tournament.Name != "" {
		d := &data{
			Applications: c.Sound.AppNames(),
			Tournament:   c.Tournament,
			Jingles:      c.Jingles.JingleStorageList(),
		}
		out, err := yaml.Marshal(d)
		if err != nil {
			c.Log.Error(err.Error())
		}

		f, err := os.Create(path.Join(c.StorageDir(), "config.yml"))

		if err != nil {
			c.Log.Error(err.Error())
		} else {
			defer f.Close()
			f.Write(out)
		}

		last, _ := os.Create(path.Join(c.AppDir(), "last.tournament"))
		last.Write([]byte(c.Tournament.Name))
		defer last.Close()

		return out
	}
	return []byte{}
}

type addableListI interface {
	AddUniq(string, LogI) (bool, error)
}

// Load - will load data from yaml file
func (c *Context) Load(input []byte) {
	c.cleanup()
	d := &data{}
	yaml.Unmarshal(input, d)

	if d.Tournament == nil {
		return
	}

	c.Tournament = d.Tournament
	c.Tournament.context = c
	c.Tournament.PlanJingles()
	for _, val := range d.Jingles {
		s, err := NewSong(val.File, c)
		if err != nil {
			c.Log.Error(err.Error())
		} else {
			c.Songs.AddUniq(s, c.Log)
			c.Jingles.AddUniq(NewJingle(val.Name, s, val.TimeBeforePoint, val.Point, c), c.Log)
		}
	}
	for _, val := range d.Applications {
		c.Log.Debug("adding application: " + val)
		c.Sound.AddUniq(val, c.Log)
	}
}

// LoadByName - will load data from file
func (c *Context) LoadByName(name string) error {
	if name == "" {
		return errors.New("Nothing to load")
	}
	f, err := os.Open(path.Join(c.AppDir(), name, "config.yml"))
	if err != nil {
		c.Log.Error("File opening error " + name + ": " + err.Error())
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.Log.Error("File read error " + name + ": " + err.Error())
		return err
	}
	c.Load(data)
	return nil
}

// SaveSong - will save uploaded song file into tournament directory
func (c *Context) SaveSong(r io.Reader, filename string) (string, error) {
	c.Log.Info(filename)
	targetFile := path.Join(c.MediaDir(), filename)
	writer, err := os.Create(targetFile)

	if err != nil {
		c.Log.Error("File upload failed: " + err.Error())
		return "", err
	}

	defer writer.Close()
	_, err = io.Copy(writer, r)

	if err != nil {
		c.Log.Error("File upload failed: " + err.Error())
		return "", err
	}
	return filename, nil
}

// RemoveSong - will remove song
func (c *Context) RemoveSong(filename string) error {
	filepath := path.Join(c.MediaDir(), filename)
	return os.Remove(filepath)
}
