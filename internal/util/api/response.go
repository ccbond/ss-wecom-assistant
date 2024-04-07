// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	resp "github.com/SyntSugar/ss-infra-go/api/response"
)

func ResponseWithOK(c *gin.Context, data any) {
	resp.ResponseWithOK(c, data)
}

func ResponseWithCreated(c *gin.Context, data any) {
	resp.ResponseWithCreated(c, data)
}

func ResponseWithSuccess(c *gin.Context, statusCode int, data any) {
	resp.ResponseWithSuccess(c, statusCode, data)
}

func ResponseErrors(c *gin.Context, errors ...any) {
	resp.ResponseWithErrors(c, http.StatusInternalServerError, 0, errors)
	c.Abort()
}

func ResponseConflict(c *gin.Context, errors ...any) {
	resp.ResponseWithErrors(c, http.StatusConflict, 0, errors)
	c.Abort()
}

func ResponseNotFound(c *gin.Context, errors ...any) {
	resp.ResponseWithErrors(c, http.StatusNotFound, 0, errors)
	c.Abort()
}

func ResponseUnprocessableEntity(c *gin.Context, errors ...any) {
	resp.ResponseWithErrors(c, http.StatusUnprocessableEntity, 0, errors)
	c.Abort()
}

func BadRequest(c *gin.Context, errors ...any) {
	resp.ResponseWithErrors(c, http.StatusBadRequest, 0, errors)
	c.Abort()
}

func ResponseForbidden(c *gin.Context, errors ...any) {
	resp.ResponseWithErrors(c, http.StatusForbidden, 0, errors)
	c.Abort()
}

func ExternalApiErrorResponse(c *gin.Context, err *APIError, errors ...any) {
	errs := []any{err.Error()}
	if len(errors) > 0 {
		errs = append(errs, errors)
	}
	resp.ResponseWithErrors(c, err.StatusCode, 0, errs)
}

func ResponseFileStream(ctx *gin.Context, filename string, fileData []byte) error {
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename = %s", filename))
	ctx.Header("Content-Type", "application/octet-stream")

	_, err := ctx.Writer.Write(fileData)
	return err
}
