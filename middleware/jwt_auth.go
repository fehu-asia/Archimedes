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

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(constant.CtxUser)
		userForm := value.(*form.UserForm)
		// 如果ctx不存在，从redis中获取
		if exists == false || userForm == nil {
			strToken := c.Request.Header.Get(constant.RequestLoginToken)
			claims, e := jwt2.VerifyAction(strToken)
			if e != nil {
				ResponseError(c, NoLoginErrorCode, errors.New(fmt.Sprintf("您需要登录！%s", e.Error())))
				c.Abort()
				return
			} else {
				tokenId := claims.TokenId
				userForm := &form.UserForm{}
				e := redis_lib.GetObject(tokenId, userForm, "")
				if e != nil {
					ResponseError(c, NoLoginErrorCode, errors.New(fmt.Sprintf("会话超时，请重新登录！")))
					c.Abort()
					return
				}
				fmt.Println("获取到当前用户:", userForm, e)
				c.Set(constant.CtxUser, userForm)
				// 当用户访问过以后，更新token
				defer redis_lib.Set(claims.TokenId, userForm, constant.UserLoginToKenExpireTime, "")
			}

		} else {
			if userForm.UserId == "" {
				ResponseError(c, NoLoginErrorCode, errors.New(fmt.Sprintf("您需要登录！")))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
