package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   *T     `json:"data"`
	Error  string `json:"error"`
}

func Ok[T any](ctx *gin.Context, data *T) {
	ctx.JSON(http.StatusOK, Response[T]{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   data,
	})
}

func OkNoData(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response[any]{
		Code:   http.StatusOK,
		Status: "SUCCESS",
	})
}

func Created[T any](ctx *gin.Context, data *T) {
	ctx.JSON(http.StatusCreated, Response[T]{
		Code:   http.StatusCreated,
		Status: "SUCCESS",
		Data:   data,
	})
}
