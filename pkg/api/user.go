package api

import (
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rovn208/ross/pkg/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"github.com/rovn208/ross/pkg/util"
)

var (
	ErrUserNameOrPasswordDoesNotCorrect = errors.New("username or password does not correct, please try again")
	ErrUserDoesNotExists                = errors.New("user does not exists")
	ErrUserAlreadyExists                = errors.New("user already exists")
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
		ID:        user.ID,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(ErrUserAlreadyExists))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type updateUserRequest struct {
	Password string `json:"password,omitempty" binding:"min=6"`
	Email    string `json:"email" binding:"email"`
	FullName string `json:"full_name,omitempty"`
	Username string `json:"username,omitempty" binding:"required,alphanum"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		Email: pgtype.Text{
			String: req.Email,
			Valid:  req.Email != "",
		},
		FullName: pgtype.Text{
			String: req.FullName,
			Valid:  req.FullName != "",
		},
		Username: req.Username,
	}

	if req.Password != "" {
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)
}

type getUserRequest struct {
	ID int64 `uri:"id,omitempty" binding:"required"`
}

func (server *Server) getUserByID(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrUserDoesNotExists))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

func (server *Server) getUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type loginRequest struct {
	Username string `json:"username,omitempty" binding:"required,min=6,alphanum"`
	Password string `json:"password,omitempty" binding:"required,min=6"`
}

type loginResponse struct {
	AccessToken           string       `json:"access_token,omitempty"`
	AccessTokenExpiredAt  time.Time    `json:"access_token_expired_at,omitempty"`
	RefreshToken          string       `json:"refresh_token,omitempty"`
	RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at,omitempty"`
	User                  userResponse `json:"user"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrUserDoesNotExists))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUserNameOrPasswordDoesNotCorrect))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loginResponse{
		User:                  newUserResponse(user),
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshTokenPayload.ExpiredAt,
	})
}
