package gin

import (
	"errors"
	"gin/SDK/html/template"
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

// Exit 中断执行
func Exit() {
	panic(errors.New("exit"))
}
