package exception

type HttpException struct {
	Message string
}

func NewHttpException(error string) *HttpException {
	return &HttpException{
		Message: "HttpException",
	}
}

func (that HttpException) Error() string {
	return that.Message
}
