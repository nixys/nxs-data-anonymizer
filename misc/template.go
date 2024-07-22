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
func TemplateExec(tpl string, d any) ([]byte, error) {

	var b bytes.Buffer

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

	err = t.Execute(&b, d)
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
