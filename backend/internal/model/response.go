package model

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string `json:"token"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// DecideResponse 决策响应
type DecideResponse struct {
	Menu    Menu   `json:"menu"`
	Message string `json:"message"`
}

// HistoryResponse 历史记录响应
type HistoryResponse struct {
	Records []DecisionRecord `json:"records"`
	Total   int64            `json:"total"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}
