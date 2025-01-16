package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Responder struct {
	C *gin.Context
}

type Response struct {
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	PageInfo *PageInfo   `json:"pageInfo,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}

type PageInfo struct {
	CurrentPage int `json:"currentPage"`
	NextPage    int `json:"nextPage"`
	PrevPage    int `json:"prevPage"`
	TotalPage   int `json:"totalPage"`
	TotalData   int `json:"totalData"`
}

func NewResponse(ctx *gin.Context) *Responder {
	return &Responder{C: ctx}
}

func (r *Responder) Success(message string, data interface{}) {
	r.C.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func (r *Responder) GetAllSuccess(message string, data interface{}, page *PageInfo) {
	r.C.JSON(http.StatusOK, Response{
		Status:   http.StatusOK,
		Message:  message,
		PageInfo: page,
		Data:     data,
	})
}

func (r *Responder) Created(message string, data interface{}) {
	r.C.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: message,
		Data:    data,
	})
}

func (r *Responder) BadRequest(message string, err interface{}) {
	r.C.JSON(http.StatusBadRequest, Response{
		Status:  http.StatusBadRequest,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) Forbidden(message string, err interface{}) {
	r.C.JSON(http.StatusForbidden, Response{
		Status:  http.StatusForbidden,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) Unauthorized(message string, err interface{}) {
	r.C.JSON(http.StatusUnauthorized, Response{
		Status:  http.StatusUnauthorized,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) NotFound(message string, err interface{}) {
	r.C.JSON(http.StatusNotFound, Response{
		Status:  http.StatusNotFound,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) InternalServerError(message string, err interface{}) {
	r.C.JSON(http.StatusInternalServerError, Response{
		Status:  http.StatusInternalServerError,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}
