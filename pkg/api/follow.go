package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"github.com/rovn208/ross/pkg/token"
	"net/http"
)

type followUserRequest struct {
	FollowingUserID int64 `json:"following_user_id" binding:"required"`
}

// FollowUser godoc
// @Summary Follow user
// @Tags follows
// @Accept json
// @Produce json
// @Param following_user_id body int64 true "123456789"
// @Success 200 {object} string "{"messsage": "follow user successfully"}"
// @Failure 400,500 {object} error "{"error": "error message"}"
// @Router /api/v1/follows/followers [post]
func (server *Server) followUser(ctx *gin.Context) {
	var req followUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.FollowUserParams{
		FollowedUserID:  authPayload.UserID,
		FollowingUserID: req.FollowingUserID,
	}

	_, err := server.store.FollowUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, messageResponse("follow user successfully"))
}

type unfollowUserRequest struct {
	FollowingUserID int64 `json:"following_user_id" binding:"required"`
}

// UnfollowUser godoc
// @Summary Unfollow user
// @Tags follows
// @Accept json
// @Produce json
// @Param following_user_id body int64 true "123456789"
// @Success 200 {object} string "{"messsage": "unfollow user successfully"}"
// @Failure 400,500 {object} error "{"error": "error message"}"
// @Router /api/v1/follows/followers [delete]
func (server *Server) unfollowUser(ctx *gin.Context) {
	var req unfollowUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UnfollowUserParams{
		FollowedUserID:  authPayload.UserID,
		FollowingUserID: req.FollowingUserID,
	}

	err := server.store.UnfollowUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, messageResponse("Unfollowed user successfully"))
}

type getListFollowRequest struct {
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

// GetListFollower godoc
// @Summary Get list follower
// @Description Get list follower
// @Tags follows
// @Produce json
// @Param limit query int32 true "20"
// @Param offset query int32 false "0"
// @Success 200 {array} listFollowResponse
// @Failure 500 {object} error "{"error": "error message"}"
// @Router /api/v1/follows/followers [get]
func (server *Server) getListFollower(ctx *gin.Context) {
	var req getListFollowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetListFollowerParams{
		FollowingUserID: authPayload.UserID,
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

// GetListFollowing godoc
// @Summary Get list following
// @Description Get list following
// @Tags follows
// @Produce json
// @Param limit query int32 true "20"
// @Param offset query int32 false "0"
// @Success 200 {array} listFollowResponse
// @Failure 500 {object} error "{"error": "error message"}"
// @Router /api/v1/follows/following [get]
func (server *Server) getListFollowing(ctx *gin.Context) {
	var req getListFollowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetListFollowingParams{
		FollowedUserID: authPayload.UserID,
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
