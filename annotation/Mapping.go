package annotation

import "gin"

type Handler func(*gin.RouterGroup, []gin.HandlerFunc, map[string]string, string) ([]string, string, error)

type mapping struct {
	Map map[string]func(...any) Handler
}

var Mapping = &mapping{
	Map: make(map[string]func(...any) Handler),
}

func (that *mapping) Register(name string, fn func(...any) Handler) {
	that.Map[name] = fn
}

func (that *mapping) Get(name string) func(...any) Handler {
	if that.Map[name] == nil {
		return nil
	}
	return that.Map[name]
}

func init() {
	Mapping.Register("Request", RequestMapping)
	Mapping.Register("Post", PostMapping)
	Mapping.Register("Get", GetMapping)
	Mapping.Register("Put", PutMapping)
	Mapping.Register("Delete", DeleteMapping)
	Mapping.Register("Patch", PatchMapping)
}
