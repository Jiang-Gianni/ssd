package parse

import (
	"fmt"
	"html/template"
	"io"
	"regexp"
)

func (s *Schema) Parse(filenames []string, entrypoint string, w io.Writer) error {
	fm := template.FuncMap{
		"regex": Regex,
	}

	tmpl, err := template.New("").Funcs(fm).ParseFiles(filenames...)
	if err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	return tmpl.ExecuteTemplate(w, entrypoint, s)
}

func Regex(pattern, value string) (bool, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	return r.MatchString(value), nil
}
