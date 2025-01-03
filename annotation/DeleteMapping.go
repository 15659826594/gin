package annotation

func DeleteMapping(args ...any) Handler {
	return RequestMapping([]string{"DELETE"})
}
