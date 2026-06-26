package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
)

type successBean struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type errorBean struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

// OKJSON writes JSON response directly (compatible with go-zero httpx.OkJsonCtx).
func OKJSON(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// OK writes an empty 200 response (compatible with go-zero httpx.Ok).
func OK(c *gin.Context) {
	c.Status(http.StatusOK)
}

// Error writes error message with 400 status (compatible with go-zero httpx.ErrorCtx).
func Error(c *gin.Context, err error) {
	c.String(http.StatusBadRequest, err.Error())
}

// Fail writes wrapped error JSON used by legacy result.Fail.
func Fail(c *gin.Context, err error) {
	log.Errorf("[API][URL:%s] - %s", c.Request.URL, err.Error())

	causeErr := errors.Cause(err)
	if gStatus, ok := status.FromError(causeErr); ok {
		c.JSON(http.StatusOK, &errorBean{Code: 500, Msg: gStatus.Message()})
		return
	}

	c.JSON(http.StatusOK, &errorBean{Code: 500, Msg: err.Error()})
}

// JSONOk writes wrapped success JSON used by legacy result.JsonOk.
func JSONOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &successBean{Code: 200, Msg: "OK", Data: data})
}
