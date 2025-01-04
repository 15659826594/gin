package app

func htmlentities(s string) string {
	return html.EscapeString(s)
}
