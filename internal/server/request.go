package server

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
}

func NewRequest() *Request {
	return &Request{
		Headers: make(map[string]string),
	}
}
