package annotation

func PostMapping(args ...any) Handler {
	return RequestMapping([]string{"POST"})
}
