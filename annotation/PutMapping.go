package annotation

func PutMapping(httpMethod string, attr map[string]string) ([]string, string) {
	return RequestMapping("PUT", attr)
}
