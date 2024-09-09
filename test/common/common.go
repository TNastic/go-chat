package main

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
	"io/ioutil"
)

func main() {

	// 测试邮件发送
	//e := email.NewEmail()
	//e.From = "Chat Craft<1497556691@qq.com>"
	//e.To = []string{"1497556691@qq.com"}
	////e.Bcc = []string{"test_bcc@example.com"}
	////e.Cc = []string{"test_cc@example.com"}
	//e.Subject = "Awesome Subject"
	//e.Text = []byte("给你的信息")
	//e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	//err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "1497556691@qq.com", "svoaswgsjiabgjah", "smtp.qq.com"),
	//	&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	//if err != nil {
	//	return
	//}

	UploadFile()

}

// UploadFile  测试文件上传
func UploadFile() {

	// ak sk 空间 key:资源路径

	// 读取本地文件
	data, _ := ioutil.ReadFile("D:/Temp/bzjdt.jpg")

	putPolicy := storage.PutPolicy{
		Scope: "lxlde",
	}
	mac := auth.New("JiqMw7t9Jx9q_u6cYzWGtS7asabP6gywL2mHnUp9", "C6SJl2ujPIfgYpvX-EEXYUy02neiWrSHBNXB2LCU")

	// 上传token
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneXinjiapo
	//bucketManager := storage.NewBucketManager(mac, &cfg)

	//fileInfo, sErr := bucketManager.Stat("lxlde", "picture")
	//if sErr == nil && fileInfo.Fsize != 0 {
	//	// 当文件在云端存在则不上传
	//	return
	//}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	dataLen := int64(len(data))
	err := formUploader.Put(context.Background(), &ret, upToken, "picture/alsl.jpg", bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
}
