package master

import (
	"bytes"
	"reflect"
	"regexp"
	"text/template"
)

// EscapedTemplate executes template without replacing undefined field
func EscapedTemplate(s string, params interface{}) (string, error) {
	r := regexp.MustCompile(`\{\{\s*(\.[a-zA-Z0-9_^\}\s]*)\s*\}\}`)
	s = r.ReplaceAllString(s, `{{ if $1 }}{{$1}}{{ else }}{{"{{$1}}"}}{{ end }}`)

	tpl, err := template.New("").Parse(s)
	if err != nil {
		return "", err
	}

	val := reflect.ValueOf(params)
	typ := val.Type()
	kind := val.Kind()

	dict := make(map[string]interface{})
	if kind == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			field := val.Field(i)
			dict[typ.Field(i).Name] = field.Interface()
		}
	} else if kind == reflect.Map {
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			dict[k.String()] = v
		}
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, dict)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
