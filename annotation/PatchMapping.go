package annotation

func PatchMapping(args ...any) Handler {
	return RequestMapping([]string{"PATCH"})
}
