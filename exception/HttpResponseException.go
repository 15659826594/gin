package exception

import "gin/response"

type HttpResponseException struct {
	Message  string
	response *response.Response
}

func NewHttpResponseException(resp *response.Response) *HttpResponseException {
	return &HttpResponseException{
		Message:  "HttpResponseException",
		response: resp,
	}
}

func (e *HttpResponseException) GetMessage() string {
	return e.Message
}

func (e *HttpResponseException) GetResponse() *response.Response {
	return e.response
}
