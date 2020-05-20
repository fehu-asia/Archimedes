package middleware

import (
	"bytes"
	"fehu/common/lib"
	"fehu/common/lib/redis_lib"
	"fehu/constant"
	"fehu/util/cryp"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func CrpyBefore(c *gin.Context) {
	request := c.Request
	data, err := c.GetRawData()
	if err != nil {
		ParamVerifyError(c)
		c.Abort()
		return
	}
	if strings.HasPrefix(request.RequestURI, "/base") ||
		strings.HasPrefix(request.RequestURI, "/login") ||
		request.RequestURI == "/ping" {
		c.Set("requestData", string(data))
		request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点
		return
	}
	// 开启了加密
	if lib.GetIntConf("base.http.crpyStatus") == 1 {
		// 签名的加密方式  = rsa(sha256(秘文),publicKey)
		signature := request.Header.Get(constant.RequestHeaderSignature)
		if signature == "" {
			SignatureError(c, "签名为空！")
			c.Abort()
			return
		}
		hex := cryp.GetSHA256HashCode(data)
		privateKey, err := redis_lib.GetString(constant.RedisRsaPrivateKey, "")
		if err != nil {
			SignatureError(c, "签名错误，没有私钥!")
			c.Abort()
			return
		}
		decrypt, err := cryp.RsaDecrypt(cryp.Base64DecodeByte(signature), privateKey)
		if err != nil {
			SignatureError(c, "签名错误！")
			c.Abort()
			return
		}
		// 私钥解密被公钥加密的签名和发送过来的签名进行对比，如果这一步出问题：
		// 1. 对称加密密钥是否一致
		// 2. 公钥匙和私钥是否能对上
		if hex != string(decrypt) {
			SignatureError(c, "参数签名不一致!")
			c.Abort()
			return
		}
		// 解包
		key := c.GetString(constant.CtxAesKey)
		if key == "" {
			EncryptionError(c, "no key!")
			c.Abort()
			return
		}
		data, err = cryp.AesEncryptByte(data, key)
		if err != nil {
			EncryptionError(c, "key error!")
			c.Abort()
			return
		}
		// 响应体是否加密，这个值不为空，就加密
		c.Set(constant.ReponseBodyCrypto, constant.ReponseBodyCrypto)

	}
	// 如果不加密，原封不动放进去，如果开启加密，把解密的数据放进去
	c.Set("requestData", string(data))
	request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点
}
func after(c *gin.Context) {

	fmt.Println("加密后")
}
func Crpy() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不是post请求
		requestUri := c.Request.RequestURI
		fmt.Println(requestUri)
		if c.Request.Method == http.MethodPost {
			CrpyBefore(c)
		}
		c.Next()
		after(c)
	}
}
