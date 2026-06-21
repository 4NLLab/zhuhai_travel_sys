package dto

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

func Success(data interface{}) *Response {
	return &Response{Code: 200, Message: "ok", Data: data}
}

func Fail(code int, msg string) *Response {
	return &Response{Code: code, Message: msg}
}

func Page(data interface{}, total int64, page, size int) *PageResponse {
	return &PageResponse{
		Code: 200, Message: "ok",
		Data: data, Total: total,
		Page: page, Size: size,
	}
}
