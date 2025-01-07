package gin

import (
	"gin/src/html/template"
)

// LoadHTMLFolder loads HTML files identified folder
// and associates the result with HTML renderer.
func (engine *Engine) LoadHTMLFolder(path string) {
	left := engine.delims.Left
	right := engine.delims.Right
	templ := template.Must(template.WrapT(template.New("").Delims(left, right).Funcs(engine.FuncMap)).ParseFolder(path))

	if IsDebugging() {
		debugPrintLoadTemplate(templ.Template)
	}

	engine.SetHTMLTemplate(templ.Template)
}
