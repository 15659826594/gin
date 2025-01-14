package annotation

func PostMapping(httpMethod string, attr map[string]string) ([]string, string) {
	return RequestMapping("POST", attr)
}
