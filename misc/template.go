package misc

import (
	"bytes"
	ttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

var null = "::NULL::"

type TemplateData struct {
	TableName string
	Values    map[string][]byte
	Variables map[string]string
}

// TemplateExec makes message from given template `tpl` and data `d`
func TemplateExec(tpl string, d *TemplateData) ([]byte, error) {

	type tplData struct {
		TableName string
		Values    map[string]string
		Variables map[string]string
	}

	var (
		b  bytes.Buffer
		td *tplData
	)

	if d != nil {
		td = &tplData{
			TableName: d.TableName,
			Values:    make(map[string]string),
			Variables: make(map[string]string),
		}

		for k, v := range d.Values {
			if v == nil {
				td.Values[k] = null
			} else {
				td.Values[k] = string(v)
			}
		}

		td.Variables = d.Variables
	}

	// See http://masterminds.github.io/sprig/ for details
	t, err := ttemplate.New("template").Funcs(func() ttemplate.FuncMap {

		// Get current sprig functions
		t := sprig.TxtFuncMap()

		// Add additional functions
		t["null"] = func() string {
			return null
		}
		t["isNull"] = func(v string) bool {
			if v == null {
				return true
			}
			return false
		}

		return t
	}()).Parse(tpl)
	if err != nil {
		return []byte{}, err
	}

	err = t.Execute(&b, td)
	if err != nil {
		return []byte{}, err
	}

	// Return empty line if buffer is nil
	if b.Bytes() == nil {
		return []byte{}, nil
	}

	// Return nil if buffer is NULL (with special key)
	if bytes.Compare(b.Bytes(), []byte(null)) == 0 {
		return nil, nil
	}

	// Return buffer content otherwise
	return b.Bytes(), nil
}
