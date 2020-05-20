package middleware

import (
	"encoding/json"
	"fehu/common/lib"
	"fehu/constant"
	"fehu/util/cryp"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type ResponseCode int

//1000以下为通用码，1000以上为用户自定义码
const (
	ErrorCode ResponseCode = iota
	SuccessCode
	UndefErrorCode
	ValidErrorCode
	InternalErrorCode

	NoLoginErrorCode        ResponseCode = -10
	NoSessionErrorCode      ResponseCode = -20
	SignatureErrorCode      ResponseCode = -88
	EncryptionErrorCode     ResponseCode = -99
	LimiterErrorCode        ResponseCode = 999
	ParamVerifyErrorCode    ResponseCode = 190
	ParamRequireErrorCode   ResponseCode = 199
	InvalidRequestErrorCode ResponseCode = 401
	CustomizeCode           ResponseCode = 1000

	GROUPALL_SAVE_FLOWERROR ResponseCode = 2001
)

type Response struct {
	ErrorCode ResponseCode `json:"code"`
	ErrorMsg  string       `json:"msg"`
	Data      interface{}  `json:"data"`
	TraceId   interface{}  `json:"traceId"`
	Stack     interface{}  `json:"stack,omitempty"`
}

func SignatureError(c *gin.Context, msg ...string) {
	finalMsg := "签名错误 !"
	if len(msg) > 0 {
		finalMsg = msg[0]
	}
	Error(c, SignatureErrorCode, finalMsg)
}
func EncryptionError(c *gin.Context, msg ...string) {
	finalMsg := "encryption err !"
	if len(msg) > 0 {
		finalMsg = msg[0]
	}
	Error(c, EncryptionErrorCode, finalMsg)
}
func ParamRequireError(c *gin.Context, msg ...string) {
	finalMsg := "param required !"
	if len(msg) > 0 {
		finalMsg = msg[0]
	}
	Error(c, ParamRequireErrorCode, finalMsg)
}
func ParamVerifyError(c *gin.Context, msg ...string) {
	finalMsg := "param verify err !"
	if len(msg) > 0 {
		finalMsg = msg[0]
	}
	Error(c, ParamRequireErrorCode, finalMsg)
}
func ErrorMsg(c *gin.Context, msg string) {
	Error(c, ErrorCode, msg)
}

func Error(c *gin.Context, code ResponseCode, msg string) {
	trace, _ := c.Get("trace")
	traceContext, _ := trace.(*lib.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}
	resp := &Response{ErrorCode: code, ErrorMsg: msg, Data: "", TraceId: traceId}

	SerializeJSON(c, resp, 200, nil)
}

func Success(c *gin.Context, datas ...interface{}) {
	var data interface{}
	if len(datas) > 0 {
		data = datas[0]
	}
	trace, _ := c.Get("trace")
	traceContext, _ := trace.(*lib.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}
	resp := &Response{ErrorCode: SuccessCode, ErrorMsg: "", Data: data, TraceId: traceId}

	SerializeJSON(c, resp, 200, nil)
}

func ResponseError(c *gin.Context, code ResponseCode, err error) {
	trace, _ := c.Get("trace")
	traceContext, _ := trace.(*lib.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}
	stack := ""
	if c.Query("is_debug") == "1" || lib.GetConfEnv() == "dev" {
		stack = strings.Replace(fmt.Sprintf("%+v", err), err.Error()+"\n", "", -1)
	}

	resp := &Response{ErrorCode: code, ErrorMsg: err.Error(), Data: "", TraceId: traceId, Stack: stack}
	SerializeJSON(c, resp, 200, err)
}
func SerializeJSON(c *gin.Context, res *Response, httpCode int, err error) {
	response, _ := json.Marshal(res)
	resJson := string(response)
	c.Set("response", resJson)

	if c.GetString(constant.ReponseBodyCrypto) != "" {
		key := c.GetString(constant.CtxAesKey)
		if key != "" {
			encrypt, _ := cryp.AesEncrypt(resJson, key)
			c.JSON(httpCode, encrypt)
		} else {
			c.JSON(httpCode, "0")
		}
	} else {
		c.JSON(httpCode, res)
	}

	if err != nil {
		c.AbortWithError(httpCode, err)
	}

}
