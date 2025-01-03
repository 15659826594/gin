package annotation

func PutMapping(args ...any) Handler {
	return RequestMapping([]string{"PUT"})
}
