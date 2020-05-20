package middleware

import (
	"errors"
	"fehu/common/lib"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"log"
	"time"
)

var (
	bucket *ratelimit.Bucket
)

//func init() {
//	conf := lib.GetIntConf("base.http.limitRequestCountPerSecond")
//	log.Printf(" [INFO] HttpServerRun:%s\n", lib.GetIntConf("base.http.limitRequestCountPerSecond"))
//	bucket = ratelimit.NewBucket(time.Duration(int(time.Second)/conf), int64(conf))
//}

func Limiter() gin.HandlerFunc {
	conf := lib.GetIntConf("base.http.limitRequestCountPerSecond")
	log.Printf(" [INFO] limiterCount:%d\n", lib.GetIntConf("base.http.limitRequestCountPerSecond"))
	bucket = ratelimit.NewBucket(time.Duration(int(time.Second)/conf), int64(conf))

	return func(c *gin.Context) {
		fmt.Println("tackAvailable", bucket.Available())
		if bucket.TakeAvailable(1) == 0 {
			ResponseError(c, LimiterErrorCode, errors.New(fmt.Sprintf("%v, limiter error !", c.ClientIP())))
			c.Abort()
			return
		}
		c.Next()
	}
}

//func LimitHandler(lmt *limiter.Limiter) gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//
//		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
//
//		if httpError != nil {
//
//			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
//
//			c.Abort()
//
//		} else {
//
//			c.Next()
//
//		}
//
//	}
//
//}
