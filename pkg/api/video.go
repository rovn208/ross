package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"net/http"
	"strconv"
	"time"
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
	ID           int64  `json:"id" binding:"required"`
	Title        string `json:"title"`
	StreamUrl    string `json:"stream_url"`
	Description  string `json:"description"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

func (server *Server) updateVideo(ctx *gin.Context) {
	var req updateVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateVideoParams{
		ID: req.ID,
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
