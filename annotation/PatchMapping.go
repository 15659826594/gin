package annotation

func PatchMapping(httpMethod string, attr map[string]string) ([]string, string) {
	return RequestMapping("PATCH", attr)
}
