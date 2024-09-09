package common

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
	"go-chat/pkg/config"
)

var (
	conf      = config.NewConfig("config/local.yml")
	ak        = conf.GetString("qiniu.access_key")
	sk        = conf.GetString("qiniu.secret_key")
	bucket    = conf.GetString("qiniu.bucket")
	key       = conf.GetString("qiniu.key_prefix")
	urlPrefix = conf.GetString("qiniu.url_prefix")
)

// FileUploadByBytes 字节方式上传文件
func FileUploadByBytes(filename string, data []byte) (string, error) {
	// token
	token := GenerateToken()
	// 文件名
	generatedFileName := GenerateFileName(filename)
	// 上传
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneXinjiapo
	// 上传器
	formUploader := storage.NewFormUploader(&cfg)
	// 返回值
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	dataLen := int64(len(data))
	err := formUploader.Put(context.Background(), &ret, token,
		fmt.Sprintf("%s%s", key, generatedFileName), bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", urlPrefix, ret.Key), nil
}

// FileUploadByStream 流式上传文件
func FileUploadByStream() {

}

func GenerateToken() string {
	// 上传策略
	putPolicy := storage.PutPolicy{
		Scope:           fmt.Sprintf("%s:%s", bucket, key),
		IsPrefixalScope: 1,
	}
	mac := auth.New(ak, sk)
	// 上传token
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

func GenerateFileName(filename string) string {
	// 获取文件后缀
	suffix := filename[len(filename)-4:]
	// 为文件名生成UUID
	uuidName := uuid.New().String()
	fileName := fmt.Sprintf("%s%s", uuidName, suffix)
	return fileName
}
