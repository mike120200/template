package cmn

import "fmt"

const (
	Success     = 0
	CommonError = -1
)

// AppError 是我们的自定义错误类型
type AppError struct {
	StatusCode int
	Message    string
}

// Error 方法让 AppError 实现了 error 接口
func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d,err: %s", e.StatusCode, e.Message)
}

func NewAppError(statusCode int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
}
