package service

import (
	"context"
	"errors"
	v1 "go-chat/api/v1"
	"go-chat/global"
	"go-chat/internal/model"
	"go-chat/internal/repository"
	"go-chat/pkg/common"
	"time"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	VerifyRegisterEmailCode(ctx context.Context, email string, code string) error
	CreateNewUser(ctx context.Context, req *v1.CheckRegisterEmailCodeRequest) (*v1.RegisterResponse, error)
	UpdateUserInfo(ctx context.Context, token string, userId int, req *v1.UpdateUserInfoRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, *model.UserBasics, error)
	EmailLoginCodeCheck(ctx context.Context, email string, code string) (string, *model.UserBasics, error)
	SendEmail(ctx context.Context, email string) error

	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
}

func NewUserService(
	service *Service,
	emailService EmailService,
	userRepo repository.UserRepository,
) UserService {
	return &userService{
		userRepo:     userRepo,
		emailService: emailService,
		Service:      service,
	}
}

type userService struct {
	userRepo     repository.UserRepository
	emailService EmailService
	*Service
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	user := &model.UserBasics{}
	user.Name = req.Name
	user.Email = req.Email
	// 获取加密密码
	encryptionPassword := req.Password
	reEncryptionPassword := req.RePassword
	// rsa解密
	password, err := common.RsaDecoder(encryptionPassword)
	if err != nil {
		return err
	}
	rePassword, err := common.RsaDecoder(reEncryptionPassword)
	if err != nil {
		return err
	}
	if password != rePassword {
		return errors.New("两次输入的密码不一致")
	}
	// 检查用户是否已经注册
	_, err = s.userRepo.FindUserByNameWithRegister(ctx, req.Name)
	if err != nil {
		return err
	}
	// 查询邮箱是否已经被注册
	_, err = s.userRepo.FindUserByEmailWithRegister(ctx, req.Email)
	if err != nil {
		return v1.ErrEmailAlreadyUse
	}
	// 发送验证码
	err = s.emailService.SendEmail(ctx, req.Email, global.Register)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) VerifyRegisterEmailCode(ctx context.Context, email string, code string) error {
	// 验证邮箱验证码
	err := s.emailService.CheckEmailCode(ctx, email, code, global.Register)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) CreateNewUser(ctx context.Context, req *v1.CheckRegisterEmailCodeRequest) (*v1.RegisterResponse, error) {
	// 验证邮箱验证码
	err := s.emailService.CheckEmailCode(ctx, req.Email, req.Code, global.Register)
	if err != nil {
		return nil, err
	}
	// 创建用户
	user := &model.UserBasics{}
	user.Email = req.Email
	user.Name = req.Name

	originPassword, err := common.RsaDecoder(req.Password)
	if err != nil {
		return nil, err
	}
	// 加密密码存储
	password := common.SaltPassWord(originPassword, "lxl")
	user.PassWord = password
	user.Salt = "lxl"
	t := time.Now()
	user.LoginTime = &t
	user.HeartBeatTime = &t
	user.LoginOutTime = &t

	// 创建用户
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// 获取用户
	user, err = s.userRepo.FindByName(ctx, user.Name)
	if err != nil {
		return nil, err
	}

	// jwt校验
	token, err := s.jwt.GenToken(user.Name, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return nil, err
	}

	// 返回数据
	return &v1.RegisterResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *userService) UpdateUserInfo(ctx context.Context, token string, userId int, req *v1.UpdateUserInfoRequest) error {
	user := &model.UserBasics{}

	// 字段校验
	if req.UserName == "" || req.Email == "" || req.Phone == "" || req.Avatar == "" {
		return v1.ErrBadRequest
	}
	user.ID = uint(userId)
	userInfo, _ := s.userRepo.FindUserInfoById(ctx, uint(userId))

	if userInfo.Avatar != req.Avatar {
		user.Avatar = req.Avatar
	}

	if userInfo.Email != req.Email {
		user.Email = req.Email
	}

	if userInfo.Motto != req.Motto {
		user.Motto = req.Motto
	}

	if userInfo.Phone != req.Phone {
		user.Phone = req.Phone
	}

	if userInfo.Name != req.UserName {
		// 检查用户是否已经注册
		_, err := s.userRepo.FindUserByNameWithRegister(ctx, req.UserName)
		if err != nil {
			return v1.ErrUserNameAlreadyUse
		}
		user.Name = req.UserName
	}

	err := s.userRepo.UpdateUserInfo(ctx, user)
	if err != nil {
		return v1.ErrUserInfoUpdateFailed
	}
	return nil
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, *model.UserBasics, error) {
	// 参数校验
	if req.Name == "" || req.Password == "" {
		return "", nil, v1.ErrBadRequest
	}
	// 获取用户信息
	user, err := s.userRepo.FindUserInfoByName(ctx, req.Name)
	if err != nil {
		return "", nil, v1.ErrUserNotFound
	}

	// 校验密码
	loginPassword, err := common.RsaDecoder(req.Password)
	if err != nil {
		return "", nil, err
	}
	if !common.CheckPassWord(loginPassword, user.Salt, user.PassWord) {
		return "", nil, v1.ErrUserPasswordError
	}

	token, err := s.jwt.GenToken(user.Name, time.Now().Add(time.Hour*24*90))

	if err != nil {
		return "", nil, v1.ErrInternalServerError
	}

	return token, user, nil
}

func (s *userService) EmailLoginCodeCheck(ctx context.Context, email string, code string) (string, *model.UserBasics, error) {

	// 获取用户信息
	user, err := s.userRepo.FindUserByEmailWithLogin(ctx, email)
	if err != nil {
		s.logger.WithContext(ctx).Error("Email not found")
		return "", nil, err
	}

	// 校验邮箱验证码
	err = s.emailService.CheckEmailCode(ctx, email, code, global.Login)
	if err != nil {
		return "", nil, err
	}

	// 生成token
	token, err := s.jwt.GenToken(user.Name, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", nil, err
	}

	return token, user, nil

}

func (s *userService) SendEmail(ctx context.Context, email string) error {

	// 验证码是否存在
	_, err := s.userRepo.FindUserByEmailWithLogin(ctx, email)
	if err != nil {
		s.logger.WithContext(ctx).Error("Email not found")
		return err
	}

	// 发送验证码
	err = s.emailService.SendEmail(ctx, email, global.Login)

	if err != nil {
		s.logger.WithContext(ctx).Error("Send email failed")
		return err
	}

	return nil

}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:   user.UserId,
		Nickname: user.Nickname,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Nickname = req.Nickname

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
