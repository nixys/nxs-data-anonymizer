package misc

import (
	"bytes"
	ttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

var null = []byte("::NULL::")

type TemplateData struct {
	TableName string
	Values    map[string][]byte
}

// TemplateExec makes message from given template `tpl` and data `d`
func TemplateExec(tpl string, d TemplateData) ([]byte, error) {

	var b bytes.Buffer

	for k, v := range d.Values {
		if v == nil {
			d.Values[k] = null
		}
	}

	// See http://masterminds.github.io/sprig/ for details
	t, err := ttemplate.New("template").Funcs(func() ttemplate.FuncMap {

		// Get current sprig functions
		t := sprig.TxtFuncMap()

		// Add additional functions
		t["null"] = func() string {
			return string(null)
		}
		t["isNull"] = func(v []byte) bool {
			if bytes.Compare(v, null) == 0 {
				return true
			}
			return false
		}

		return t
	}()).Parse(tpl)
	if err != nil {
		return []byte{}, err
	}

	err = t.Execute(&b, d)
	if err != nil {
		return []byte{}, err
	}

	// Return empty line if buffer is nil
	if b.Bytes() == nil {
		return []byte{}, nil
	}

	// Return nil if buffer is NULL (with special key)
	if bytes.Compare(b.Bytes(), null) == 0 {
		return nil, nil
	}

	// Return buffer content otherwise
	return b.Bytes(), nil
}
