package gooh

type MemoryContext struct {
	data map[string]interface{}
}

func (c *MemoryContext) Get(k string) (interface{}, error) {
	return c.data[k], nil
}

func (c *MemoryContext) Set(k string, d interface{}) error {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[k] = d
	return nil
}

func (c *MemoryContext) Exists(k string) (bool, error) {
	_, ok := c.data[k]
	return ok, nil
}
