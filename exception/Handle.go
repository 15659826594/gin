package exception

import "gin"

func Handle(c *gin.Context) {
	defer func() {
		rec := recover()
		//捕获自定义异常
		if _, ok := rec.(Exception); ok {
			switch exception := rec.(type) {
			case *HttpResponseException:
				exception.GetResponse().Send(c)
				c.Abort()
			}
		} else if rec != nil {
			panic(rec)
		}
	}()
	c.Next()
}
