package api

import (
	"net/http"
	db "simpleauth/db/sqlc"
	"simpleauth/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createPostRequest struct {
	Title           string `json:"title" binding:"required"`
	PostDescription string `json:"post_description" binding:"required"`
}

func (server *Server) CreatePost(ctx *gin.Context) {
	authPayload := ctx.MustGet(authoriaztionPayloadKey).(*token.Payload)
	var req createPostRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreatePostParams{
		ID:              uuid.New().String(),
		Title:           req.Title,
		PostDescription: req.PostDescription,
		AuthorName:      authPayload.Username,
		PostDate:        time.Now(),
	}
	post, err := server.store.CreatePost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, post)
}

func (server *Server) ListPosts(ctx *gin.Context) {
	posts, err := server.store.ListPosts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, posts)
}
