package controller

import (
	"fehu/common/lib/redis_lib"
	"fehu/constant"
	"fehu/middleware"
	"fehu/model/form"
	"fehu/model/param"
	"fehu/util/cryp"
	"fehu/util/jwt"
	"fehu/util/map_builder"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type LoginController struct {
}

func NoLoginRequiredRegister(router *gin.RouterGroup) {
	l := LoginController{}
	router.POST("/login", l.login)
}

//登录，获取jwt token
func (l *LoginController) login(c *gin.Context) {
	var loginParam param.LoginParam
	err := c.BindJSON(&loginParam)
	if err != nil {
		fmt.Println(err)
	}
	if true {
		claims := &jwt.JWTClaims{
			//UserID:   1,
			//Username: loginParam.Username,
			//Password: loginParam.Password,
			TokenId: cryp.GenUUID(),
		}
		claims.IssuedAt = time.Now().Unix()
		claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(jwt.ExpireTime)).Unix()
		fmt.Println(claims.ExpiresAt)
		singedToken, err := jwt.GenToken(claims)
		if err != nil {
			middleware.ErrorMsg(c, "口令生成失败！")
			return
		}
		//map[string]interface{}{"token": singedToken}
		// 将用户信息放到redis中
		userForm := form.UserForm{
			Username: loginParam.Username,
			UserId:   cryp.GenUUID(),
			AesKey:   c.GetString(constant.CtxAesKey),
		}
		err = redis_lib.Set(claims.TokenId, userForm, constant.UnLoginToKenExpireTime, "")
		if err != nil {
			middleware.ErrorMsg(c, "登录失败！")
			return
		}
		middleware.Success(c, map_builder.BuilderMap(
			"token", singedToken,
			"test", 1234))
	}
}
