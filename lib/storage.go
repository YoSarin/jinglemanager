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
	Songs        []string
	Applications []string
	Jingles      []*Jingle
}

// Save - will save data to yaml file
func (c *Context) Save() []byte {
	if c.Tournament.Name != "" {
		d := &data{
			Songs:        c.Songs.FileNames(),
			Applications: c.Sound.AppNames(),
			Tournament:   c.Tournament,
			Jingles:      c.Jingles.JingleList(),
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
	for _, val := range d.Songs {
		c.Log.Debug("adding song: " + val)
        s, err := NewSong(val, c.Log)
        if err != nil {
            c.Log.Error(err.Error())
        } else {
            c.Songs.AddUniq(s, c.Log)
        }
	}
	for _, val := range d.Applications {
		c.Log.Debug("adding application: " + val)
		c.Sound.AddUniq(val, c.Log)
	}
    if d.Tournament != nil {
        c.Tournament = d.Tournament
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
	targetFile := path.Join(c.StorageDir(), "media", filename)
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
	return targetFile, nil
}
