package response

type Result struct {
	Code   int               `json:"code"`
	Msg    string            `json:"msg"`
	Time   int64             `json:"time"`
	Data   any               `json:"data"`
	Type   string            `json:"-"`
	Header map[string]string `json:"-"`
}
