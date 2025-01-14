package annotation

func GetMapping(httpMethod string, attr map[string]string) ([]string, string) {
	return RequestMapping("GET", attr)
}
