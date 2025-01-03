package gin

import (
	"fmt"
	"gin/lib/gotable"
	"github.com/gin-gonic/gin"
)

func DebugPrintTable(title []string, fn func(func([]string))) {
	if gin.IsDebugging() {
		var content [][]string
		content = make([][]string, 0)
		fn(func(row []string) {
			content = append(content, row)
		})
		table, _ := gotable.Create(title...)
		for i := 0; i < len(content); i++ {
			_ = table.AddRow(content[i])
		}
		fmt.Print(table)
	} else {
		fn(func(row []string) {})
	}
}
