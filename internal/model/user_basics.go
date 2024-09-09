package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID       uint           `json:"id" gorm:"primaryKey"`
	CreateAt time.Time      `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt time.Time      `json:"update_at" gorm:"autoUpdateTime"`
	DeleteAt gorm.DeletedAt `json:"delete_at" gorm:"index"`
}

type UserBasics struct {
	Model
	Name          string     `json:"name" gorm:"name"`
	PassWord      string     `json:"pass_word" gorm:"pass_word"`
	Avatar        string     `json:"avatar" gorm:"avatar"`
	Gender        string     `json:"gender" gorm:"gender"`
	Phone         string     `json:"phone" gorm:"phone"`
	Email         string     `json:"email" gorm:"email"`
	Motto         string     `json:"motto" gorm:"motto"`
	Identity      string     `json:"identity" gorm:"identity"`
	ClientIp      string     `json:"client_ip" gorm:"client_ip"`
	ClientPort    string     `json:"client_port" gorm:"client_port"`
	Salt          string     `json:"salt" gorm:"salt"`
	LoginTime     *time.Time `json:"login_time" gorm:"login_time"`
	HeartBeatTime *time.Time `json:"heart_beat_time" gorm:"heart_beat_time"`
	LoginOutTime  *time.Time `json:"login_out_time" gorm:"login_out_time"`
	IsLoginOut    int8       `json:"is_login_out" gorm:"is_login_out"`
	DeviceInfo    string     `json:"device_info" gorm:"device_info"`
}

func (*UserBasics) TableName() string {
	return "user_basics"
}
