package annotation

type Handler func(string, map[string]string) ([]string, string)

var RouteMapping = map[string]Handler{}

func Register(name string, fn Handler) {
	RouteMapping[name] = fn
}

func Get(name string) Handler {
	if fn, ok := RouteMapping[name]; ok {
		return fn
	}
	return nil
}

func init() {
	Register("Request", RequestMapping)
	Register("Post", PostMapping)
	Register("Get", GetMapping)
	Register("Put", PutMapping)
	Register("Delete", DeleteMapping)
	Register("Patch", PatchMapping)
}
