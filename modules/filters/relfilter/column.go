package relfilter

import "fmt"

type columns struct {
	cc []*column
	m  map[string]*column
}

type column struct {
	n string
	t columnTypes
}

type columnTypes struct {
	raw    string
	groups [][]string
	r      *ColumnRuleOpts
	env    []string
}

func columnsInit() columns {
	return columns{
		cc: []*column{},
		m:  make(map[string]*column),
	}
}

func (c *columns) add(name string, rt string, pts [][]string, r *ColumnRuleOpts) {

	env := []string{fmt.Sprintf("%s=%s", envVarColumnTypeRAW, rt)}

	if pts != nil {
		for i, g := range pts {
			for j, sg := range g {

				if j == 0 {
					env = append(
						env,
						fmt.Sprintf("%s%d=%s", envVarColumnTypeGroupPrefix, i, sg),
					)
				} else {
					env = append(
						env,
						fmt.Sprintf("%s%d_%d=%s", envVarColumnTypeGroupPrefix, i, j-1, sg),
					)
				}
			}
		}
	}

	v := column{
		n: name,
		t: columnTypes{
			raw:    rt,
			groups: pts,
			r:      r,
			env:    env,
		},
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
