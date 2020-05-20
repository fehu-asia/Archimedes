package controller

import (
	"fehu/constant"
	"fehu/middleware"
	"github.com/gin-gonic/gin"
)

type ApiController struct {
}

func ApiRegister(router *gin.RouterGroup) {
	api := ApiController{}
	router.POST("/needLogin", api.NeedLogin)
}

//登录，获取jwt token
func (a *ApiController) NeedLogin(c *gin.Context) {
	value, _ := c.Get(constant.CtxUser)
	middleware.Success(c, value)
}
