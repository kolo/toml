package toto

import (
	"os"
)

type Conf struct {
	values map[string]interface{}
}

func (c *Conf) String(key string) string {
	if v, ok := c.values[key]; ok {
		return v.(string)
	} else {
		return ""
	}
}

func (c *Conf) Int(key string) int64 {
	if v, ok := c.values[key]; ok {
		return v.(int64)
	} else {
		var zero int64
		return zero
	}
}

func (c *Conf) Bool(key string) bool {
	if v, ok := c.values[key]; ok {
		return v.(bool)
	} else {
		var boolean bool
		return boolean
	}
}

func (c *Conf) Slice(key string) []interface{} {
	if v, ok := c.values[key]; ok {
		return v.([]interface{})
	} else {
		return nil
	}
}

func newConf() *Conf {
	return &Conf{
		values: make(map[string]interface{}),
	}
}

func Parse(path string) (conf *Conf, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf = newConf()
	err = parse(f, conf)
	if err != nil {
		return nil, err
	}

	return
}
