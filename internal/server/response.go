package server

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func NewResponse() *Response {
	return &Response{
		Headers: make(map[string]string),
	}
}
