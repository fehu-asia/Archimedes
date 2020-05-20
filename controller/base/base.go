package base

import (
	"encoding/json"
	"fehu/common/lib/redis_lib"
	"fehu/constant"
	"fehu/middleware"
	"fehu/model/param"
	"fehu/util/cryp"
	"github.com/gin-gonic/gin"
	"strings"
)

type BaseController struct {
}

func BaseRegister(router *gin.RouterGroup) {
	base := BaseController{}
	router.GET("/img", base.Img)
	router.POST("/sayHi", base.SayHi)
}

/**

 */
func (b *BaseController) Img(c *gin.Context) {
	hex := c.Query("hex")
	if hex != "" {
		// 判断是否和原有hex一致,不要重复发送公钥匙
		if strings.Trim(hex, " ") == "1234" {
			c.Data(201, "image/png", []byte{})
			return
		}
	}
	// hex 为空，或者不想等代表客户端公钥已经过期，重新发送公钥
	publickKey, e := redis_lib.GetString(constant.RedisRsaPublicKey, "")
	if e != nil || publickKey == "" {
		privateKeyBytes, publicKeyBytes := cryp.GenerateRsaKey(2048)
		e = redis_lib.Set(constant.RedisRsaPublicKey, string(publicKeyBytes), 3600*24, "")
		if e != nil {
			c.Data(500, "image/png", []byte{})
			return
		}
		e = redis_lib.Set(constant.RedisRsaPrivateKey, string(privateKeyBytes), 3600*24, "")
		if e != nil {
			c.Data(500, "image/png", []byte{})
			return
		}
		c.Header("content-disposition", `attachment; filename=app01.png`)
		s := cryp.Base64EncodeByteCount(publicKeyBytes, 5)
		c.Data(200, "image/png", []byte(s))
	} else {
		c.Header("content-disposition", `attachment; filename=app01.png`)
		s := cryp.Base64EncodeByteCount([]byte(publickKey), 5)
		c.Data(200, "image/png", []byte(s))
	}

}

func (b *BaseController) SayHi(c *gin.Context) {
	request := c.Request
	data, err := c.GetRawData()
	if err != nil {
		middleware.EncryptionError(c, "[1]解密失败！")
		return
	}
	privateKey, err := redis_lib.GetString(constant.RedisRsaPrivateKey, "")
	if err != nil {
		middleware.EncryptionError(c, "[2]解密失败！")
		return
	}
	jsonbytes, err := cryp.RsaDecrypt(data, privateKey)
	if err != nil {
		middleware.EncryptionError(c, "[3]解密失败！")
		return
	}
	// 解密成功，做签名校验
	signature := request.Header.Get("signature")
	sha256Hex, err := cryp.RsaDecrypt([]byte(signature), privateKey)
	if err != nil {
		middleware.SignatureError(c, "[1]签名错误！")
		return
	}
	if code := cryp.GetSHA256HashCode(jsonbytes); code != string(sha256Hex) {
		middleware.SignatureError(c, "[2]签名错误！")
		return
	}

	// 转换json做业务逻辑规则校验
	var param param.BaseParam
	err = json.Unmarshal(jsonbytes, &param)
	if err != nil {
		middleware.ParamVerifyError(c, "参数解析失败!")
		return
	}
	// 1.版本号
	if param.Version == "" {
		middleware.ParamRequireError(c)
		return
	}
	// 生成对称加密密钥，发送给前端
	uuId := cryp.GenUUID()
	//claims := &jwt.JWTClaims{
	//	TokenId: cryp.GenUUID(),
	//}
	//claims.IssuedAt = time.Now().Unix()
	//claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(jwt.ExpireTime)).Unix()
	//singedToken, err := jwt.GenToken(claims)
	//if err != nil {
	//	middleware.ErrorMsg(c, "口令生成失败！")
	//	return
	//}
	//// uuid 放到 redis中
	//userForm := form.UserForm{
	//	AesKey: c.GetString(constant.CtxAesKey),
	//}
	//err = redis_lib.Set(claims.TokenId, userForm, constant.UserLoginToKenExpireTime, "")
	//if err != nil {
	//	middleware.ErrorMsg(c, "设置缓存失败！")
	//	return
	//}
	bytes, err := cryp.PriKeyENCTYPT([]byte(uuId), privateKey)
	// 发送对称加密密钥
	middleware.Success(c, map[string]string{"key": cryp.Base64EncodeByte(bytes)})

}
