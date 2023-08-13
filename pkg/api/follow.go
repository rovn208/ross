package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"net/http"
)

type followUserRequest struct {
	UserID          int64 `json:"user_id" binding:"required"`
	FollowingUserID int64 `json:"following_user_id" binding:"required"`
}

func (server *Server) followUser(ctx *gin.Context) {
	var req followUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: Check current userid
	arg := db.FollowUserParams{
		FollowedUserID:  req.UserID,
		FollowingUserID: req.FollowingUserID,
	}

	_, err := server.store.FollowUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, messageResponse("follow user successfully"))
}

type getListFollowRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	Limit  int32 `json:"limit" binding:"required,min=1"`
	Offset int32 `json:"offset"` // Offset 0 if it's not provided
}

type listFollowResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

func (server *Server) newListFollowerResponse(ctx *gin.Context, follows []db.Follow) ([]listFollowResponse, error) {
	followers := make([]listFollowResponse, 0)
	for _, follow := range follows {
		f, err := server.store.GetUser(ctx, follow.FollowedUserID)
		if err != nil {
			return followers, err
		}
		followers = append(followers, listFollowResponse{
			UserID:   f.ID,
			Username: f.Username,
			FullName: f.FullName,
		})
	}

	return followers, nil
}

func (server *Server) getListFollower(ctx *gin.Context) {
	var req getListFollowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: Check current userid
	arg := db.GetListFollowerParams{
		FollowingUserID: req.UserID,
		Limit:           req.Limit,
		Offset:          req.Offset,
	}

	follows, err := server.store.GetListFollower(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	followers, err := server.newListFollowerResponse(ctx, follows)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, followers)
}

func (server *Server) newListFollowingResponse(ctx *gin.Context, follows []db.Follow) ([]listFollowResponse, error) {
	followers := make([]listFollowResponse, 0)
	for _, follow := range follows {
		f, err := server.store.GetUser(ctx, follow.FollowingUserID)
		if err != nil {
			return followers, err
		}
		followers = append(followers, listFollowResponse{
			UserID:   f.ID,
			Username: f.Username,
			FullName: f.FullName,
		})
	}

	return followers, nil
}

func (server *Server) getListFollowing(ctx *gin.Context) {
	var req getListFollowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: Check current userid
	arg := db.GetListFollowingParams{
		FollowedUserID: req.UserID,
		Limit:          req.Limit,
		Offset:         req.Offset,
	}

	follows, err := server.store.GetListFollowing(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	followers, err := server.newListFollowingResponse(ctx, follows)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, followers)
}

type unfollowUserRequest struct {
	UserID          int64 `json:"user_id" binding:"required"`
	FollowingUserID int64 `json:"following_user_id" binding:"required"`
}

func (server *Server) unfollowUser(ctx *gin.Context) {
	var req unfollowUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: Check current userid
	arg := db.UnfollowUserParams{
		FollowedUserID:  req.UserID,
		FollowingUserID: req.FollowingUserID,
	}

	err := server.store.UnfollowUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, messageResponse("Unfollowed user successfully"))
}
