package middleware

import (
	"errors"
	"fehu/common/lib/redis_lib"
	"fehu/constant"
	"fehu/model/form"
	jwt2 "fehu/util/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 拦截所有请求函数，只要前端有token就去尝试读取并更新token
		strToken := c.Request.Header.Get(constant.RequestLoginToken)
		if strToken == "" {
			c.Next()
			return
		}
		claims, e := jwt2.VerifyAction(strToken)
		if e != nil {
			ResponseError(c, NoSessionErrorCode, errors.New(fmt.Sprintf("no session！%s", e.Error())))
			c.Abort()
			return
		} else {
			tokenId := claims.TokenId
			userForm := &form.UserForm{}
			e := redis_lib.GetObject(tokenId, userForm, "")
			if e != nil {
				ResponseError(c, NoSessionErrorCode, errors.New(fmt.Sprintf("会话超时，请重新发起会话！")))
				c.Abort()
				return
			}
			c.Set(constant.CtxUser, userForm)
			c.Set(constant.CtxAesKey, tokenId)
			// 当用户访问过以后，更新token
			defer redis_lib.Set(claims.TokenId, userForm, constant.UserLoginToKenExpireTime, "")
		}

	}
}
