package middleware

import (
	"bytes"
	"fehu/common/lib"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xinliangnote/go-util/time"
	"log"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

/**
写往通道中的日志
*/
type ChanLog struct {
	log map[string]interface{}
	c   *gin.Context
	t   *lib.TraceContext
}

var accessChannel = make(chan *ChanLog, 100)

//func RequestLog() gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//		SetUp(c)
//		//defer RequestOutLog(c)
//		//c.Next()
//	}
//}

func RequestLog() gin.HandlerFunc {
	go handleAccessChannel()

	return func(c *gin.Context) {
		traceContext := lib.NewTrace()
		if traceId := c.Request.Header.Get("com-header-rid"); traceId != "" {
			traceContext.TraceId = traceId
		}
		if spanId := c.Request.Header.Get("com-header-spanid"); spanId != "" {
			traceContext.SpanId = spanId
		}
		c.Set("trace", traceContext)
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// 开始时间
		startTime := time.GetCurrentMilliUnix()

		// 处理请求
		c.Next()

		responseBody, _ := c.Get("response")

		// 结束时间
		endTime := time.GetCurrentMilliUnix()

		// 日志格式
		accessLogMap := make(map[string]interface{})

		accessLogMap["request_time"] = startTime
		accessLogMap["request_method"] = c.Request.Method
		accessLogMap["request_uri"] = c.Request.RequestURI
		accessLogMap["request_proto"] = c.Request.Proto
		accessLogMap["request_ua"] = c.Request.UserAgent()
		accessLogMap["request_referer"] = c.Request.Referer()
		data, _ := c.Get("requestData")
		accessLogMap["request_raw_data"] = data
		accessLogMap["request_client_ip"] = c.ClientIP()

		accessLogMap["response_time"] = endTime
		accessLogMap["response"] = responseBody

		accessLogMap["cost_time"] = fmt.Sprintf("%vms", endTime-startTime)

		//accessLogJson, _ := jsonUtil.Encode(accessLogMap)
		accessChannel <- &ChanLog{log: accessLogMap, c: c, t: traceContext}
	}
}

func handleAccessChannel() {
	//if f, err := os.OpenFile(config.AppAccessLogName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
	//	log.Println(err)
	//} else {
	//	for accessLog := range accessChannel {
	//		_, _ = f.WriteString(accessLog + "\n")
	//	}
	//}
	//return
	for chanLog := range accessChannel {
		log.Println(chanLog.log)
		lib.Log.TagInfo(chanLog.t, "[request log]", chanLog.log)
	}

}
