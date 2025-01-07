package gin

import "errors"

// Exit 中断执行
func Exit() {
	panic(errors.New("exit"))
}

// RecoveryExit 正常中断后续请求
func RecoveryExit() HandlerFunc {
	return func(c *Context) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if err := rec.(error); err.Error() == "exit" {
				c.Abort()
			} else {
				panic(rec)
			}
		}()
		c.Next()
	}
}
