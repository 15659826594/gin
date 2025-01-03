package annotation

func GetMapping(args ...any) Handler {
	return RequestMapping([]string{"GET"})
}
