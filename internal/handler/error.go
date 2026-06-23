package handler

import (
	"errors"
	"go-loan-api/internal/apperror"
	"go-loan-api/internal/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ToErrorResponse(ctx *gin.Context, err error) {
	if appError, ok := errors.AsType[*apperror.AppError](err); ok {
		ctx.JSON(appError.Code, dto.Response[any]{
			Code:   appError.Code,
			Status: appError.Status,
			Error:  appError.Message,
		})
		return
	}

	log.Printf("Unexpected error: %v", err)
	ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL_SERVER_ERROR",
		Error:  "Something Went wrong",
	})
}
