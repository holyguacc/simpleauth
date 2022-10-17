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

func (server *Server) sendVerificationCode(ctx *gin.Context) {
	username := ctx.Query("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("no username entered")))
		return
	}
	userData, err := server.store.GetUser(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("this user does not exist")))
		return
	}
	verificationCode := uuid.New().String()
	hashedVerificationCode, err := util.HashPassword(verificationCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//set Hashed Verification Key in database
	verArg := db.SetVerifyCodeParams{
		VerifyKey: hashedVerificationCode,
		Username:  userData.Username,
	}
	err = server.store.SetVerifyCode(ctx, verArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("could not set verification key")))
	}

	//send verification code via email
	emailtype := &verificationEmail{}
	server.sendEmail(userData, verificationCode, emailtype) //set verification status (verified) to false
	verficationArg := db.SetVeriftyStatuesParams{
		IsVerified: sql.NullBool{Bool: false, Valid: true},
		VerefiedOn: time.Time{},
	}
	server.store.SetVeriftyStatues(ctx, verficationArg)
}

func (server *Server) verifyUser(ctx *gin.Context) {

	username := ctx.Query("username")
	requestCode := ctx.Query("code")

	userData, err := server.store.GetVerifyKey(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("cannot retrive user"))
		return
	}
	if err := util.CheckPassword(requestCode, userData.VerifyKey); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("code you entered does not exist")))
		return
	}
	arg := db.SetVeriftyStatuesParams{
		Username:   userData.Username,
		IsVerified: sql.NullBool{Bool: true, Valid: true},
		VerefiedOn: time.Now(),
	}
	if err := server.store.SetVeriftyStatues(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error changing verification satus %s", err)))
		return
	}
	response := fmt.Sprintf("Hi %s , Verification confirmed.", userData.Username)
	ctx.JSON(http.StatusOK, gin.H{"detail": response})

}
