package v1

import "go-chat/internal/model"

type RegisterRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required" example:"1234@gmail.com"`
	Password   string `json:"password" binding:"required" example:"123456"`
	RePassword string `json:"rePassword" binding:"required" example:"123456"`
}

type CheckRegisterEmailCodeRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
	Code     string `json:"code" binding:"required" example:"123456"`
}

type RegisterResponse struct {
	Response
	Token string            `json:"token"`
	User  *model.UserBasics `json:"user"`
}

type UpdateUserInfoRequest struct {
	UserName   string `json:"userName"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	Motto      string `json:"motto"`
	ClientIp   string `json:"clientIp"`
	ClientPort string `json:"clientPort"`
}

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginResponseData struct {
	Token    string            `json:"token"`
	UserInfo *model.UserBasics `json:"user"`
}

type EmailLoginCheckRequest struct {
	Email string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Code  string `json:"code" binding:"required" example:"123456"`
}

type EmailLoginRequest struct {
	Email string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}

type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}
