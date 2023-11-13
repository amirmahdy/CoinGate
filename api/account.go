package api

import (
	"db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Coin     string `json:"coin" binding:"required,coin"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	coin, err := server.store.GetCoin(ctx, req.Coin)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Username: user.Username,
		Coin:     coin.Name,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
