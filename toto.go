package toto

import (
	"os"
)

type Conf struct {
	values map[string]interface{}
}

func (c *Conf) String(key string) string {
	return c.values[key].(string)
}

func Parse(path string) (conf *Conf, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf = new(Conf)
	err = parse(f, conf)
	if err != nil {
		return nil, err
	}

	return
}
