package gin

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

// LoadHTMLFolder loads HTML files identified folder
// and associates the result with HTML renderer.
func (engine *Engine) LoadHTMLFolder(path string, rename func(name string) string) {
	left := engine.delims.Left
	right := engine.delims.Right
	templ := template.Must(ParseFolder(template.New("").Delims(left, right).Funcs(engine.FuncMap), path, rename))

	if IsDebugging() {
		debugPrintLoadTemplate(templ)
	}

	engine.SetHTMLTemplate(templ)
}

// ParseFolder 遍历目录查找 html和tpl文件
func ParseFolder(t *template.Template, folder string, renameFunc func(path string) string) (*template.Template, error) {
	tplFiles := map[string]string{}
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if slices.Contains([]string{".html", ".tpl", ".tmpl"}, filepath.Ext(path)) {
			if renameFunc == nil {
				tplFiles[filepath.ToSlash(path)] = path
			} else {
				tplFiles[renameFunc(path)] = path
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return parseFilesMap(t, readFileOS, tplFiles)
}

func parseFilesMap(t *template.Template, readFile func(string) (string, []byte, error), filenames map[string]string) (*template.Template, error) {
	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}
	for name, filename := range filenames {
		_, b, err := readFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)

		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
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
