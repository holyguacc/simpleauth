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

func (server *Server) SendResetPassword(ctx *gin.Context) {
	username := ctx.Query("username")
	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("could not retrieve user data")))
		return
	}
	resetKey := uuid.New().String()
	hashedResetKey, err := util.HashPassword(resetKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to generate reset Key")))
		return
	}
	if err := server.sendEmail(user, resetKey, &ResetEmail{}); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to send reset Email")))
		return
	}
	//store reset key and Put is_reset on False
	arg := db.SetResetKeyAndStateParams{
		Username: user.Username,
		ResetKey: sql.NullString{
			String: hashedResetKey,
			Valid:  true,
		},
		IsReset: sql.NullBool{Bool: false, Valid: true},
	}
	if err := server.store.SetResetKeyAndState(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("cannot set reset key"))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "reset code is emailed to you"})
}

type resetAndChangeRequest struct {
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
}

func (server *Server) ResetAndChangePassword(ctx *gin.Context) {
	resetKey := ctx.Query("code")
	var req resetAndChangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("bad request")))
		return
	}
	user, err := server.store.GetResetKey(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("no user by the username %s was found", req.Username)))
		return
	}
	if err := util.CheckPassword(resetKey, user.ResetKey.String); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("the reset key `%s` you provieded is false", resetKey)))
		return
	}
	newHashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("could not hash new password")))
		return
	}
	arg := db.UpdatePasswordParams{
		Username:          user.Username,
		HashedPassword:    newHashedPassword,
		PasswordChangedAt: time.Now(),
	}
	if err := server.store.UpdatePassword(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("could not update the old password")))
	}
	ctx.JSON(http.StatusOK, gin.H{"detail": "password changed successfuly"})
}
