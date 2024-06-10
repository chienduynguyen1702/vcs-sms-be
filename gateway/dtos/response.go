package dtos

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type Pagination struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

func ErrorResponse(message string) Response {
	return Response{
		Success: false,
		Message: message,
		Data:    nil,
	}
}
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

const (
	LIMIT_DEFAULT = 10
	PAGE_DEFAULT  = 1
)

func PaginationResponse(total int64, data interface{}) Pagination {
	return Pagination{
		Total: total,
		Data:  data,
	}
}
