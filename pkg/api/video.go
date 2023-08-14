package api

import (
	"errors"
	"fmt"
	"github.com/rovn208/ross/pkg/token"
	"github.com/rovn208/ross/pkg/util"
	"github.com/rovn208/ross/pkg/youtube"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"github.com/rovn208/ross/pkg/ffmpeg"
)

var (
	ErrVideoDoesNotExists = errors.New("video does not exists")
	ErrVideoAlreadyExists = errors.New("video already exists")
)

type createVideoRequest struct {
	Title        string `json:"title,omitempty" binding:"required, min=6"`
	StreamUrl    string `json:"stream_url,omitempty" binding:"required"`
	Description  string `json:"description,omitempty" `
	ThumbnailUrl string `json:"thumbnail_url,omitempty" `
	CreatedBy    int64  `json:"created_by,omitempty" binding:"required"`
}

type videoResponse struct {
	ID           int64       `json:"id"`
	Title        string      `json:"title"`
	StreamUrl    string      `json:"stream_url"`
	Description  pgtype.Text `json:"description"`
	ThumbnailUrl pgtype.Text `json:"thumbnail_url"`
	CreatedBy    int64       `json:"created_by"`
	CreatedAt    time.Time   `json:"created_at"`
}

func newVideoResponse(video db.Video) videoResponse {
	return videoResponse{
		ID:           video.ID,
		Title:        video.Title,
		StreamUrl:    video.StreamUrl,
		Description:  video.Description,
		ThumbnailUrl: video.ThumbnailUrl,
		CreatedBy:    video.CreatedBy,
		CreatedAt:    video.CreatedAt,
	}
}

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

type addYoutubeVideoRequest struct {
	URL string `json:"url,omitempty" binding:"required"`
}

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

	ctx.JSON(http.StatusInternalServerError, messageResponse("creating video via youtube video successfully"))
}

func toListVideoResponse(videos []db.Video) []videoResponse {
	res := make([]videoResponse, len(videos))
	for i, video := range videos {
		res[i] = newVideoResponse(video)
	}
	return res
}

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

func (server *Server) uploadVideo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "To be implemented")
}
