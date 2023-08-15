package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rovn208/ross/pkg/token"
	"github.com/rovn208/ross/pkg/util"
	"github.com/rovn208/ross/pkg/youtube"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"github.com/rovn208/ross/pkg/ffmpeg"
)

var (
	ErrVideoDoesNotExists       = errors.New("video does not exists")
	ErrVideoAlreadyExists       = errors.New("video already exists")
	ErrUnsupportedFileExtension = errors.New("unsupported file extension")
)

type createVideoRequest struct {
	Title        string `json:"title,omitempty" binding:"required, min=6"`
	StreamUrl    string `json:"stream_url,omitempty" binding:"required"`
	Description  string `json:"description,omitempty" `
	ThumbnailUrl string `json:"thumbnail_url,omitempty" `
	CreatedBy    int64  `json:"created_by,omitempty" binding:"required"`
}

type videoResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	StreamUrl    string    `json:"stream_url"`
	Description  string    `json:"description"`
	ThumbnailUrl string    `json:"thumbnail_url"`
	CreatedBy    int64     `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

func newVideoResponse(video db.Video) videoResponse {
	return videoResponse{
		ID:           video.ID,
		Title:        video.Title,
		StreamUrl:    video.StreamUrl,
		Description:  video.Description.String,
		ThumbnailUrl: video.ThumbnailUrl.String,
		CreatedBy:    video.CreatedBy,
		CreatedAt:    video.CreatedAt,
	}
}

// CreateVideo godoc
// @Summary Create new video
// @Description Create new video
// @Tags video
// @Accept json
// @Produce json
// @Param title body string true "Video title" minlength(6)
// @Param stream_url body string true "foldername/video.mp4"
// @Param description body string false "Video description"
// @Param thumbnail_url body string false "https://i.ytimg.com/vi/-uFQzcY7YHc/hqdefault.jpg?sqp=-oaymwEbCKgBEF5IVfKriqkDDggBFQAAiEIYAXABwAEG&rs=AOn4CLBdAOc5E4H_G09C5wqorhYRsUwQrQ"
// @Param created_by body int64 true "123451"
// @Success 200 {object} videoResponse
// @Failure 400,500 {object} error "{"error": "error message"}"
// @Router /api/v1/videos [post]
func (server *Server) createVideo(ctx *gin.Context) {
	var req createVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateVideoParams{
		Title:     req.Title,
		StreamUrl: req.StreamUrl,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Description != "",
		},
		ThumbnailUrl: pgtype.Text{
			String: req.ThumbnailUrl,
			Valid:  req.ThumbnailUrl != "",
		},
		CreatedBy: req.CreatedBy,
	}

	video, err := server.store.CreateVideo(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, newVideoResponse(video))
}

type videoIDUriRequest struct {
	ID string `uri:"id" binding:"required"`
}

// DeleteVideo godoc
// @Summary Delete video
// @Description Delete video
// @Tags video
// @Produce json
// @Param id path int64 true "12345"
// @Success 200 {object} string "deleted video successfully"
// @Failure 400,500 {object} error "{"error": "error message"}"
// @Router /api/v1/videos/{id} [delete]
func (server *Server) deleteVideo(ctx *gin.Context) {
	var req videoIDUriRequest
	id, err := bindAndGetIdUri(req, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = server.store.DeleteVideo(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("deleted video successfully"))

}

// GetVideo godoc
// @Summary Get video
// @Description Get video
// @Tags video
// @Produce json
// @Param id path int64 true "12345"
// @Success 200 {object} videoResponse
// @Failure 400,500 {object} error "{"error": "error message"}"
// @Router /api/v1/videos/{id} [get]
func (server *Server) getVideo(ctx *gin.Context) {
	var req videoIDUriRequest
	id, err := bindAndGetIdUri(req, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	video, err := server.store.GetVideo(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newVideoResponse(video))
}

type updateVideoRequest struct {
	Title        string `json:"title"`
	StreamUrl    string `json:"stream_url"`
	Description  string `json:"description"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

// UpdateVideo godoc
// @Summary Update video
// @Description Update video
// @Tags video
// @Accept json
// @Produce json
// @Param title body string false "Video title" minlength(6)
// @Param stream_url body string false "foldername/video.mp4"
// @Param description body string false "Video description"
// @Param thumbnail_url body string false "https://i.ytimg.com/vi/-uFQzcY7YHc/hqdefault.jpg?sqp=-oaymwEbCKgBEF5IVfKriqkDDggBFQAAiEIYAXABwAEG&rs=AOn4CLBdAOc5E4H_G09C5wqorhYRsUwQrQ"
// @Success 200 {object} videoResponse
// @Failure 400,500 {object} error "{"error": "error message"}"
// @Router /api/v1/videos/{id} [put]
func (server *Server) updateVideo(ctx *gin.Context) {
	var uri videoIDUriRequest
	id, err := bindAndGetIdUri(uri, ctx)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrVideoDoesNotExists))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateVideoRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateVideoParams{
		ID: id,
		Title: pgtype.Text{
			String: req.Title,
			Valid:  req.Title != "",
		},
		StreamUrl: pgtype.Text{
			String: req.StreamUrl,
			Valid:  req.StreamUrl != "",
		},
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Description != "",
		},
		ThumbnailUrl: pgtype.Text{
			String: req.ThumbnailUrl,
			Valid:  req.ThumbnailUrl != "",
		},
	}

	video, err := server.store.UpdateVideo(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newVideoResponse(video))
}

func bindAndGetIdUri(req videoIDUriRequest, ctx *gin.Context) (int64, error) {
	if err := ctx.ShouldBindUri(&req); err != nil {
		return 0, err
	}
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func toListVideoResponse(videos []db.Video) []videoResponse {
	res := make([]videoResponse, len(videos))
	for i, video := range videos {
		res[i] = newVideoResponse(video)
	}
	return res
}

// GetListVideo godoc
// @Summary Get list video
// @Description Get list video
// @Tags video
// @Produce json
// @Success 200 {array} videoResponse
// @Failure 500 {object} error "{"error": "error message"}"
// @Router /api/v1/videos [get]
func (server *Server) getListVideo(ctx *gin.Context) {
	arg := db.GetListVideoParams{
		Limit:  20,
		Offset: 0,
	}
	videos, err := server.store.GetListVideo(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, toListVideoResponse(videos))
}

type addYoutubeVideoRequest struct {
	URL string `json:"url,omitempty" binding:"required"`
}

// AddYoutubeVideo godoc
// @Summary Add video via youtube video url
// @Description Add video via youtube video url
// @Tags video
// @Accept json
// @Produce json
// @Param url body string true "https://www.youtube.com/watch?v=-uFQzcY7YHc"
// @Success 200 {object} string "created video successfully"
// @Failure 400,500 {object} error "{"error": "error message"}"
// @Router /api/v1/videos/youtube [post]
func (server *Server) addYoutubeVideo(ctx *gin.Context) {
	var req addYoutubeVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	videoID, err := server.ytClient.GetVideoID(req.URL)
	if err != nil {
		util.Logger.Error("Error when getting youtube video id", "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if util.IsFileAlreadyExists(youtube.GetStreamFile(server.config, videoID)) {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrVideoAlreadyExists))
		return
	}

	ytVideo, err := server.ytClient.DownloadVideo(req.URL)
	if err != nil {
		util.Logger.Error("Error when loading youtube video", "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = ffmpeg.ToHLSFormat(ctx, ytVideo.Name())
	if err != nil {
		util.Logger.Error("Error when converting hls", "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = os.Remove(ytVideo.Name())
	if err != nil {
		util.Logger.Error("Error when removing file", "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.CreateVideo(ctx, db.CreateVideoParams{
		StreamUrl: fmt.Sprintf("%s/%s.m3u8", videoID, videoID),
		Title:     ytVideo.Video.Title,
		Description: pgtype.Text{
			String: ytVideo.Video.Description,
			Valid:  ytVideo.Video.Description != "",
		},
		ThumbnailUrl: pgtype.Text{
			String: ytVideo.Video.Thumbnails[0].URL,
			Valid:  ytVideo.Video.Thumbnails[0].URL != "",
		},
		CreatedBy: authPayload.UserID,
	})
	if err != nil {
		util.Logger.Error("Error when saving video metadata into database", "video", ytVideo.Video, "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("created video successfully"))
}

type uploadVideoRequest struct {
	Title        string `form:"title,omitempty" binding:"required,min=8"`
	Description  string `form:"description,omitempty"`
	ThumbnailUrl string `form:"thumbnail_url,omitempty"`
}

// UploadVideo godoc
// @Summary Add video via form uploading
// @Description Add video via form uploading
// @Tags video
// @Accept multipart/form-data
// @Produce json
// @Param title body string true "Video title" minlength(6)
// @Param description body string false "Video description"
// @Param thumbnail_url body string false "https://i.ytimg.com/vi/-uFQzcY7YHc/hqdefault.jpg?sqp=-oaymwEbCKgBEF5IVfKriqkDDggBFQAAiEIYAXABwAEG&rs=AOn4CLBdAOc5E4H_G09C5wqorhYRsUwQrQ"
// @Success 200 {object} string "created video successfully"
// @Failure 400,500 {object} error "{"error": "error message"}"
func (server *Server) uploadVideo(ctx *gin.Context) {
	var form uploadVideoRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	thumbnailUrl := ctx.PostForm("thumbnail_url")
	file, err := ctx.FormFile("file")
	id := uuid.New()
	if err != nil {
		util.Logger.Error("get form file error", "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	util.Logger.Info("Creating youtube file", "videoID", id)
	dir := fmt.Sprintf("%s/%s", server.config.VideoDir, id)
	if err = util.CreateDirectory(dir); err != nil {
		util.Logger.Error("error when creating directory", "dir", dir, "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	filenameParts := strings.Split(file.Filename, ".")
	videoExtension := filenameParts[1]
	if !util.IsSupportedExtensions(videoExtension) {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrUnsupportedFileExtension))
		return
	}
	fileName := fmt.Sprintf("%s/%s.%s", dir, id, videoExtension)
	ctx.SaveUploadedFile(file, fileName)

	err = ffmpeg.ToHLSFormat(ctx, fileName)
	if err != nil {
		util.Logger.Error("Error when converting hls", "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = os.Remove(fileName)
	if err != nil {
		util.Logger.Error("Error when removing file", "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.CreateVideo(ctx, db.CreateVideoParams{
		StreamUrl: fmt.Sprintf("%s/%s.m3u8", id, id),
		Title:     title,
		Description: pgtype.Text{
			String: description,
			Valid:  description != "",
		},
		ThumbnailUrl: pgtype.Text{
			String: thumbnailUrl,
			Valid:  thumbnailUrl != "",
		},
		CreatedBy: authPayload.UserID,
	})
	if err != nil {
		util.Logger.Error("Error when saving video metadata into database", "video", fileName, "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("created video successfully"))
}
