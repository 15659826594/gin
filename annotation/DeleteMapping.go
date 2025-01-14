package annotation

func DeleteMapping(httpMethod string, attr map[string]string) ([]string, string) {
	return RequestMapping("DELETE", attr)
}
