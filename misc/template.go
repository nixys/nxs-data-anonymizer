package misc

import (
	"bytes"
	ttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	null = "::NULL::"
	drop = "::DROP::"
)

type TemplateData struct {
	TableName string
	Values    map[string][]byte
	Variables map[string]string
}

// TemplateExec makes message from given template `tpl` and data `d`
func TemplateExec(tpl string, d any) ([]byte, bool, error) {

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
		t["drop"] = func() string {
			return drop
		}

		return t
	}()).Parse(tpl)
	if err != nil {
		return []byte{}, false, err
	}

	err = t.Execute(&b, d)
	if err != nil {
		return []byte{}, false, err
	}

	// Return empty line if buffer is nil
	if b.Bytes() == nil {
		return []byte{}, false, nil
	}

	// Return nil if buffer is NULL (with special key)
	if bytes.Compare(b.Bytes(), []byte(null)) == 0 {
		return nil, false, nil
	}

	// Return `drop` value if buffer is DROP (with special key)
	if bytes.Compare(b.Bytes(), []byte(drop)) == 0 {
		return nil, true, nil
	}

	// Return buffer content otherwise
	return b.Bytes(), false, nil
}
