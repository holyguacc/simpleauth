package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "simpleauth/db/sqlc"
	"simpleauth/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//User request : data that user sends in
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"` // accepting only alphanumerical
	Password string `json:"password" binding:"required,min=6"`    //min length for password is 6
	FullName string `json:"full_name" binding:"required"`         // only english letters
	Email    string `json:"email" binding:"required,email"`
}

// what we show back to users
type createUserResponse struct {
	Username          string    `json:"username"`
	Fullname          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
type getUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type userRespose struct {
	Token string             `json:"token"`
	Data  createUserResponse `json:"data"`
}

func initUserResponse(user *db.User) createUserResponse {
	return createUserResponse{
		Username:          user.Username,
		Fullname:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	//validate if user request is correct
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// create a verfication code
	verificationCode := uuid.New().String()
	hashedVerificationCode, err := util.HashPassword(verificationCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//set Hashed Verification Key in database
	verArg := db.SetVerifyCodeParams{
		VerifyKey: hashedVerificationCode,
		Username:  user.Username,
	}
	err = server.store.SetVerifyCode(ctx, verArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//send verification code via email
	if err := server.sendEmail(user, verificationCode, &verificationEmail{}); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} //set verification status (verified) to false
	verficationArg := db.SetVeriftyStatuesParams{
		IsVerified: sql.NullBool{Bool: false, Valid: true},
		VerefiedOn: time.Time{},
	}
	server.store.SetVeriftyStatues(ctx, verficationArg)
	rsp := createUserResponse{
		Username:          user.Username,
		Fullname:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.PasswordChangedAt,
	}
	ctx.JSON(http.StatusOK, rsp)

}

func (server *Server) LoginUser(ctx *gin.Context) {

	var req getUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	status, err := server.store.GetverificationStatus(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("error retrieving user data"))
		return
	}
	if !status.Bool {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("user still not verified \n check your email")))
		return
	}
	token, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, userRespose{
		Token: token,
		Data:  initUserResponse(&user),
	})
}
