package template

import (
	"bytes"
	"text/template"
)

// Interpolate returns "Hello, John!" from source=`Hello, {{ .user }}!` and replacements=map[string]any{"user": "John"}
func Interpolate(source string, replacements map[string]any) (string, error) {
	// see https://pkg.go.dev/text/template
	parsed, err := template.New(`pkg/template/interpolate`).
		Option(`missingkey=error`).
		Parse(source)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = parsed.Execute(buf, replacements)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func MustInterpolate(source string, replacements map[string]any) string {
	res, err := Interpolate(source, replacements)
	if err != nil {
		panic(err)
	}
	return res
}
