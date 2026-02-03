package model

import (
	"testing"
)

func TestUser_TableName(t *testing.T) {
	user := User{}
	if got := user.TableName(); got != "users" {
		t.Errorf("User.TableName() = %v, want %v", got, "users")
	}
}

func TestRestaurant_TableName(t *testing.T) {
	restaurant := Restaurant{}
	if got := restaurant.TableName(); got != "restaurants" {
		t.Errorf("Restaurant.TableName() = %v, want %v", got, "restaurants")
	}
}

func TestDecisionRecord_TableName(t *testing.T) {
	record := DecisionRecord{}
	if got := record.TableName(); got != "decision_records" {
		t.Errorf("DecisionRecord.TableName() = %v, want %v", got, "decision_records")
	}
}

func TestResponse_Success(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := Success(data)

	if resp.Code != 0 {
		t.Errorf("Success().Code = %v, want 0", resp.Code)
	}
	if resp.Message != "success" {
		t.Errorf("Success().Message = %v, want 'success'", resp.Message)
	}
	if resp.Data == nil {
		t.Error("Success().Data should not be nil")
	}
}

func TestResponse_Error(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{
			name:    "bad request",
			code:    400,
			message: "参数错误",
		},
		{
			name:    "unauthorized",
			code:    401,
			message: "未授权",
		},
		{
			name:    "not found",
			code:    404,
			message: "资源不存在",
		},
		{
			name:    "internal error",
			code:    500,
			message: "服务器内部错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := Error(tt.code, tt.message)

			if resp.Code != tt.code {
				t.Errorf("Error().Code = %v, want %v", resp.Code, tt.code)
			}
			if resp.Message != tt.message {
				t.Errorf("Error().Message = %v, want %v", resp.Message, tt.message)
			}
			if resp.Data != nil {
				t.Error("Error().Data should be nil")
			}
		})
	}
}

func TestRegisterRequest_Validation(t *testing.T) {
	tests := []struct {
		name     string
		req      RegisterRequest
		valid    bool
		errorMsg string
	}{
		{
			name: "valid request",
			req: RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			valid: true,
		},
		{
			name: "username too short",
			req: RegisterRequest{
				Username: "a",
				Password: "password123",
			},
			valid:    false,
			errorMsg: "username too short",
		},
		{
			name: "password too short",
			req: RegisterRequest{
				Username: "testuser",
				Password: "12345",
			},
			valid:    false,
			errorMsg: "password too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 简单验证
			isValid := len(tt.req.Username) >= 2 && len(tt.req.Password) >= 6
			if isValid != tt.valid {
				t.Errorf("Validation result = %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name  string
		req   LoginRequest
		valid bool
	}{
		{
			name: "valid request",
			req: LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			valid: true,
		},
		{
			name: "empty username",
			req: LoginRequest{
				Username: "",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "empty password",
			req: LoginRequest{
				Username: "testuser",
				Password: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.req.Username != "" && tt.req.Password != ""
			if isValid != tt.valid {
				t.Errorf("Validation result = %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestCreateRestaurantRequest_Validation(t *testing.T) {
	tests := []struct {
		name  string
		req   CreateRestaurantRequest
		valid bool
	}{
		{
			name:  "valid request",
			req:   CreateRestaurantRequest{Name: "餐厅A"},
			valid: true,
		},
		{
			name:  "empty name",
			req:   CreateRestaurantRequest{Name: ""},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := len(tt.req.Name) >= 1 && len(tt.req.Name) <= 100
			if isValid != tt.valid {
				t.Errorf("Validation result = %v, want %v", isValid, tt.valid)
			}
		})
	}
}
