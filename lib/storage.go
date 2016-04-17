package lib

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

type data struct {
	Songs        []string
	Applications []string
}

// Save - will save data to yaml file
func Save(c *Context) []byte {
	d := &data{
		Songs:        c.Songs.FileNames(),
		Applications: c.Sound.AppNames(),
	}
	out, err := yaml.Marshal(d)
	if err != nil {
		c.Log.Error(err.Error())
	}
	f, err := os.Create("last.yml")
	if err != nil {
		c.Log.Error(err.Error())
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
func Load(c *Context, input []byte) {
	c.cleanup()
	d := &data{}
	yaml.Unmarshal(input, d)
	for _, val := range d.Songs {
		c.Log.Debug("adding song: " + val)
		c.Songs.AddUniq(val, c.Log)
	}
	for _, val := range d.Applications {
		c.Log.Debug("adding application: " + val)
		c.Sound.AddUniq(val, c.Log)
	}
}

// LoadFromFile - will load data from file
func LoadFromFile(c *Context, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		c.Log.Error("File opening error " + filename + ": " + err.Error())
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.Log.Error("File read error " + filename + ": " + err.Error())
		return err
	}
	Load(c, data)
	return nil
}
