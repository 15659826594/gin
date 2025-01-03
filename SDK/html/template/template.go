package template

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Template struct {
	*template.Template
}

func WrapT(t *template.Template) *Template {
	return &Template{t}
}

// New allocates a new HTML template with the given name.
func New(name string) *Template {
	tmpl := &Template{
		Template: template.New(name),
	}
	return tmpl
}

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFolder(pattern string) (*Template, error) {
	return parseFolder(nil, pattern)
}

func (t *Template) ParseFolder(pattern string) (*Template, error) {
	return parseFolder(t, pattern)
}

func parseFolder(t *Template, folder string) (*Template, error) {
	tplFiles := map[string]string{}
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			name := strings.TrimPrefix(path, folder)
			name = strings.TrimPrefix(strings.ReplaceAll(name, "\\", "/"), "/")
			tplFiles[name] = path
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return parseFilesMap(t, readFileOS, tplFiles)
}

func parseFilesMap(t *Template, readFile func(string) (string, []byte, error), filenames map[string]string) (*Template, error) {
	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}
	for name, filename := range filenames {
		_, b, err := readFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)

		var tmpl *Template
		if t == nil {
			t = New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = &Template{t.New(name)}
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func readFileOS(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}
