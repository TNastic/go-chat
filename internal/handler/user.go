package handler

import (
	"github.com/gin-gonic/gin"
	"go-chat/api/v1"
	"go-chat/internal/service"
	"net/http"
	"strconv"
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

// Register godoc
// @Summary 注册
// @Schemes
// @Description 用户注册发送验证码
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
// @Success 200 {object} v1.Response
// @Router /register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	// 处理登录
	var req v1.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, true)

}

// CheckRegisterEmailCode godoc
// @Summary 验证注册邮箱验证码
// @Schemes
// @Description 验证注册时的邮箱验证码，并注册新用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.CheckRegisterEmailCodeRequest true "params"
// @Success 200 {object} v1.Response
// @Router /register/check [post]
func (h *UserHandler) CheckRegisterEmailCode(ctx *gin.Context) {
	var req v1.CheckRegisterEmailCodeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 验证验证码是否正确
	res, err := h.userService.CreateNewUser(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	// 正确，开始注册
	v1.HandleSuccess(ctx, res)
}

// UserInfoUpdate 更新用户信息
func (h *UserHandler) UserInfoUpdate(ctx *gin.Context) {
	var req v1.UpdateUserInfoRequest
	// 从header中取出token和id
	token := ctx.GetHeader("Authorization")
	userId, _ := strconv.Atoi(ctx.GetHeader("userId"))

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.UpdateUserInfo(ctx, token, userId, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 返回token和userInfo
	accessToken, userInfo, err := h.userService.Login(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.LoginResponseData{
		Token:    accessToken,
		UserInfo: userInfo,
	})

}

func (h *UserHandler) EmailLoginCodeCheck(ctx *gin.Context) {
	var req v1.EmailLoginCheckRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	token, user, err := h.userService.EmailLoginCodeCheck(ctx, req.Email, req.Code)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, &v1.LoginResponseData{
		Token:    token,
		UserInfo: user,
	})

}

func (h *UserHandler) EmailLogin(ctx *gin.Context) {

	var req v1.EmailLoginRequest

	// 获取email
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 邮箱验证码登录
	if err := h.userService.SendEmail(ctx, req.Email); err != nil {
		h.logger.WithContext(ctx).Error("邮件发送失败")
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, true)

}

// GetProfile godoc
// @Summary 获取用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetProfileResponse
// @Router /user [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, user)
}

// UpdateProfile godoc
// @Summary 修改用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateProfileRequest true "params"
// @Success 200 {object} v1.Response
// @Router /user [put]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)

	var req v1.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.UpdateProfile(ctx, userId, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}
