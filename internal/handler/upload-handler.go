package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chat-socio/backend/configuration"
	"github.com/chat-socio/backend/internal/presenter"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/chat-socio/backend/pkg/storage"
	"github.com/cloudwego/hertz/pkg/app"
)

type UploadHandler struct {
	Storage storage.ObjectStorage
	Obs     *observability.Observability
}

func NewUploadHandler(storage storage.ObjectStorage, obs *observability.Observability) *UploadHandler {
	return &UploadHandler{Storage: storage, Obs: obs}
}

func (h *UploadHandler) UploadFile(ctx context.Context, c *app.RequestContext) {
	ctx, span := h.Obs.StartSpan(ctx, "UploadHandler.UploadFile")
	defer span()

	bucketName := c.FormValue("bucket_name")
	objectName := c.FormValue("object_name")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[*presenter.UploadResponse]{
			Message: err.Error(),
		})
		return
	}

	reader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[*presenter.UploadResponse]{
			Message: err.Error(),
		})
		return
	}

	err = h.Storage.PutObject(ctx, string(bucketName), string(objectName), reader, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[*presenter.UploadResponse]{
			Message: err.Error(),
		})
		return
	}

	uri, err := h.Storage.GetObjectURI(ctx, string(bucketName), string(objectName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[*presenter.UploadResponse]{
			Message: err.Error(),
		})
		return
	}

	url := fmt.Sprintf("%s%s", configuration.ConfigInstance.Minio.PublicEndpoint, uri)

	c.JSON(http.StatusOK, presenter.BaseResponse[*presenter.UploadResponse]{
		Data:    &presenter.UploadResponse{URL: url},
		Message: "File uploaded successfully",
	})
}
