package gooh

type MemoryContext struct {
	data map[string]interface{}
}

func (c *MemoryContext) init() {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
}

func (c *MemoryContext) Get(k string) (interface{}, error) {
	c.init()
	return c.data[k], nil
}

func (c *MemoryContext) Set(k string, d interface{}) error {
	c.init()
	c.data[k] = d
	return nil
}

func (c *MemoryContext) Exists(k string) (bool, error) {
	c.init()
	_, ok := c.data[k]
	return ok, nil
}
