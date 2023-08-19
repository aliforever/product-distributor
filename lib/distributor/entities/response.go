package entities

type Response struct {
	Ok      bool        `json:"ok"`
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponse(ok bool, message *string, data interface{}) *Response {
	return &Response{
		Ok:      ok,
		Message: message,
		Data:    data,
	}
}

func NewOkResponse(data interface{}) *Response {
	return NewResponse(true, nil, data)
}

func NewErrorResponse(message string) *Response {
	return NewResponse(false, &message, nil)
}
