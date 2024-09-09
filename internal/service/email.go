package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	v1 "go-chat/api/v1"
	"go-chat/internal/repository"
	"math/rand"
	"net/smtp"
	"time"
)

type EmailService interface {
	SendEmail(ctx context.Context, email string, emailType string) error
	CheckEmailCode(ctx context.Context, email string, code string, emailType string) error
}

func NewEmailService(
	service *Service,
	emailRepository repository.EmailRepository,
) EmailService {
	return &emailService{
		Service:         service,
		emailRepository: emailRepository,
	}
}

type emailService struct {
	*Service
	emailRepository repository.EmailRepository
}

func (s *emailService) SendEmail(ctx context.Context, emailDetail string, emailType string) error {
	// 获取六位随机验证码
	emailCode := GenerateRandomCode(6)
	// 生成验证码key emil + "类型"
	emailKey := GenerateEmailKey(emailType, emailDetail)
	// 存储验证码到redis
	err := s.emailRepository.SaveVerifyCode(ctx, emailKey, emailCode)
	if err != nil {
		return err
	}
	// 发送邮件
	e := email.NewEmail()
	e.From = "Lxl chat <1497556691@qq.com>"
	e.To = []string{emailDetail}
	e.Subject = "登陆验证码"
	e.HTML = []byte(fmt.Sprintf("<h1>注册验证码: %s</h1>", emailCode))
	err = e.SendWithTLS("smtp.qq.com:465",
		smtp.PlainAuth("", "1497556691@qq.com", "svoaswgsjiabgjah", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		return v1.ErrSendEmailFailed
	}
	return nil
}

func (s *emailService) CheckEmailCode(ctx context.Context, email string, code string, emailType string) error {
	emailKey := GenerateEmailKey(emailType, email)
	// 从redis中获取验证码
	emailCode, err := s.emailRepository.GetEmailCodeByKey(ctx, emailKey)
	if err != nil {
		return err
	}
	// 比较验证码
	if emailCode != code {
		return v1.ErrEmailCodeError
	}

	return nil
}

func GenerateRandomCode(codeLen int) string {
	s := "1234567890"
	code := ""
	// import random seed, or that random code will always be the same one
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < codeLen; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

func GenerateEmailKey(emailType string, email string) string {
	return fmt.Sprintf("%s-code-%s", email, emailType)
}
