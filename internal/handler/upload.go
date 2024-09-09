package handler

import (
	"github.com/gin-gonic/gin"
	v1 "go-chat/api/v1"
	"go-chat/pkg/common"
	"io"
	"net/http"
)

type UploadHandler struct {
	*Handler
}

func NewUploadHandler(
	handler *Handler,
) *UploadHandler {
	return &UploadHandler{
		Handler: handler,
	}
}

func (h *UploadHandler) Upload(ctx *gin.Context) {
	// 获取表单信息
	file, err := ctx.FormFile("file")

	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrFileUploadError, nil)
		return
	}
	defer fileContent.Close()

	fileBytes, err := io.ReadAll(fileContent)

	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrFileUploadError, nil)
		return
	}

	imageUrl, err := common.FileUploadByBytes(file.Filename, fileBytes)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrFileUploadError, nil)
		return
	}

	v1.HandleSuccess(ctx, imageUrl)

}
