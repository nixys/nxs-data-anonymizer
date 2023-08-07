package relfilter

type columns struct {
	cc []*column
	m  map[string]*column
}

type column struct {
	n string
	t ColumnType
}

type ColumnType string

const (
	ColumnTypeNone   ColumnType = "none"
	ColumnTypeString ColumnType = "string"
	ColumnTypeBinary ColumnType = "binary"
	ColumnTypeInt    ColumnType = "int"
)

func (c ColumnType) String() string {
	return string(c)
}

func columnsInit() columns {
	return columns{
		cc: []*column{},
		m:  make(map[string]*column),
	}
}

func (c *columns) add(name string, t ColumnType) {

	v := column{
		n: name,
		t: t,
	}

	c.cc = append(c.cc, &v)
	c.m[name] = &v
}

func (c *columns) typeGetByIndex(index int) ColumnType {
	if index >= len(c.cc) {
		return ColumnTypeNone
	}
	return c.cc[index].t
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
