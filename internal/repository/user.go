package repository

import (
	"context"
	"errors"
	v1 "go-chat/api/v1"
	"go-chat/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.UserBasics) error
	Update(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	FindUserByNameWithRegister(ctx context.Context, name string) (*model.UserBasics, error)
	FindUserByEmailWithRegister(ctx context.Context, name string) (*model.UserBasics, error)
	FindUserByEmailWithLogin(ctx context.Context, email string) (*model.UserBasics, error)
	FindByName(ctx context.Context, name string) (*model.UserBasics, error)
	FindUserInfoByName(ctx context.Context, name string) (*model.UserBasics, error)
	UpdateUserInfo(ctx context.Context, userInfo *model.UserBasics) error
	FindUserInfoById(ctx context.Context, id uint) (*model.UserBasics, error)
}

func NewUserRepository(
	r *Repository,
) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

type userRepository struct {
	*Repository
}

func (r *userRepository) Create(ctx context.Context, user *model.UserBasics) error {
	if err := r.DB(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userId string) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("user_id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindUserByNameWithRegister 查询用户名是否已存在
func (r *userRepository) FindUserByNameWithRegister(ctx context.Context, name string) (*model.UserBasics, error) {
	user := model.UserBasics{}
	if tx := r.DB(ctx).Where(" name= ?", name).First(&user); tx.RowsAffected == 1 {
		// 有记录
		r.logger.WithContext(ctx).Error("user already exist")
		return nil, errors.New("user already exist")
	}
	return &user, nil
}

func (r *userRepository) FindUserByEmailWithRegister(ctx context.Context, email string) (*model.UserBasics, error) {
	user := model.UserBasics{}
	if tx := r.DB(ctx).Where(" email= ?", email).First(&user); tx.RowsAffected == 1 {
		r.logger.WithContext(ctx).Error("email already exist")
		return nil, errors.New("email already exist")
	}
	return &user, nil
}

func (r *userRepository) FindByName(ctx context.Context, name string) (*model.UserBasics, error) {
	var user model.UserBasics
	if err := r.DB(ctx).Where(" name= ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserInfoByName(ctx context.Context, name string) (*model.UserBasics, error) {
	var user model.UserBasics
	if err := r.DB(ctx).Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUserInfo(ctx context.Context, userInfo *model.UserBasics) error {
	// 只更新user中的非空值
	if err := r.DB(ctx).Where(" id= ?", userInfo.ID).Updates(model.UserBasics{
		Avatar:   userInfo.Avatar,
		ClientIp: userInfo.ClientIp,
		Email:    userInfo.Email,
		Motto:    userInfo.Motto,
		Phone:    userInfo.Phone,
		Name:     userInfo.Name,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindUserInfoById(ctx context.Context, id uint) (*model.UserBasics, error) {
	var user model.UserBasics
	if err := r.DB(ctx).Where(" id= ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserByEmailWithLogin(ctx context.Context, email string) (*model.UserBasics, error) {
	var user model.UserBasics
	if err := r.DB(ctx).Where(" email= ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrUserEmailNotFound
		}
		return nil, err
	}
	return &user, nil
}
