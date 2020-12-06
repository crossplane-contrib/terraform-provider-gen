package provider

import (
	"bytes"
	"io"
	"os"

	"sigs.k8s.io/yaml"
)

type Config struct {
	Name               string   `json:"name"`
	BaseCRDVersion     string   `json:"base-crd-version"`
	PackagePath        string   `json:"package-path"`
	ExcludeResources   []string `json:"exclude-resources"`
	ExcludeResourceMap map[string]bool
}

func (c Config) IsExcluded(resourceName string) bool {
	_, ok := c.ExcludeResourceMap[resourceName]
	return ok
}

func ConfigFromFile(path string) (Config, error) {
	c := Config{}
	fh, err := os.Open(path)
	defer fh.Close()
	if err != nil {
		return c, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, fh)
	err = yaml.Unmarshal(buf.Bytes(), &c)
	if err != nil {
		return c, err
	}
	c.ExcludeResourceMap = make(map[string]bool)
	for _, er := range c.ExcludeResources {
		c.ExcludeResourceMap[er] = true
	}
	return c, nil
}
