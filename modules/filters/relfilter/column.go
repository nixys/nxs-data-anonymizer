package relfilter

type columns struct {
	cc []*column
	m  map[string]*column
}

type column struct {
	n       string
	rawType string
}

func columnsInit() columns {
	return columns{
		cc: []*column{},
		m:  make(map[string]*column),
	}
}

func (c *columns) add(name string, rt string) {

	v := column{
		n:       name,
		rawType: rt,
	}

	c.cc = append(c.cc, &v)
	c.m[name] = &v
}

func (c *columns) getNameByIndex(index int) string {
	if index >= len(c.cc) {
		return ""
	}
	return c.cc[index].n
}

func (c *columns) delByName(name string) {

	// Get current column element
	v := c.m[name]

	// Delete element from map
	delete(c.m, name)

	// Delete element from slice
	for k, e := range c.cc {
		if e == v {
			c.cc = append(c.cc[:k], c.cc[k+1:]...)
			break
		}
	}
}
