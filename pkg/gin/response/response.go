package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

const (
	ginApiResponseKey = "gin_api_response_key"
)

const (
	SuccessCode int32 = iota
	InvalidParamCode
	InternalServer
	UnauthorizedCode
)

type Response struct {
	Code    int32  `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data"`
}

func InternalServerError(c *gin.Context, err error) {
	ginJSON(c, http.StatusInternalServerError, ParseError(err))
}

func InvalidParamError(c *gin.Context, message string) {
	ginJSON(c, http.StatusBadRequest, &Response{
		Code:    InvalidParamCode,
		Message: "invalid params, " + message,
	})
}

func Unauthorized(c *gin.Context) {
	ginJSON(c, http.StatusUnauthorized, &Response{
		Code:    UnauthorizedCode,
		Message: "unauthorized",
	})
}

func Success(c *gin.Context, data any) {
	ginJSON(c, http.StatusOK, &Response{
		Code:    SuccessCode,
		Message: "success",
		Data:    data,
	})
}

func ginJSON(c *gin.Context, code int, resp *Response) {
	c.Set(ginApiResponseKey, resp)
	c.JSON(code, resp)
}

func ParseError(err error) *Response {
	code := InternalServer
	msg := "internal server error"

	if grpcErr, ok := status.FromError(err); ok {
		code = int32(grpcErr.Code())
		msg = grpcErr.Message()
	}

	return &Response{
		Code:    code,
		Message: msg,
	}
}

func GetApiResp(c *gin.Context) *Response {
	val, ok := c.Get(ginApiResponseKey)
	if !ok {
		return nil
	}
	resp, _ := val.(*Response)
	return resp
}
