package repository

import (
	"context"
	"errors"
	"time"
)

type EmailRepository interface {
	SaveVerifyCode(ctx context.Context, key string, codeValue string) error
	GetEmailCodeByKey(ctx context.Context, key string) (string, error)
}

func NewEmailRepository(
	repository *Repository,
) EmailRepository {
	return &emailRepository{
		Repository: repository,
	}
}

type emailRepository struct {
	*Repository
}

func (r *emailRepository) SaveVerifyCode(ctx context.Context, key string, codeValue string) error {
	// setNx 插入redis，插入失败说明验证码未过期
	success, err := r.rdb.SetNX(ctx, key, codeValue, time.Duration(5)*time.Minute).Result()

	if err != nil {
		return err
	}

	if !success {
		return errors.New("验证码未过期")
	}
	return nil
}

func (r *emailRepository) GetEmailCodeByKey(ctx context.Context, key string) (string, error) {
	success, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return success, nil
}
