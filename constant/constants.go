package constant

// 加解密使用的http请求头里的签名
const RequestHeaderSignature = "signature"

// 确认登录里面的token
const RequestLoginToken = "token"

// ctx里面的用户
const CtxUser = "ctxUser"

// 生成jwt和用户密码时使用的盐
const Salt = "@#1$%^&*3(_)opg#2$%^&*(5HG~|?/-+%;6."

// redis 中rsa公钥匙
const RedisRsaPublicKey = "rsa_private_key"

// redis 中rsa私钥
const RedisRsaPrivateKey = "rsa_public_key"

// ctx中的对称加密密钥对
const CtxAesKey = "ctxAesKey"

// 用户token超时时间 （20分钟）
const UserLoginToKenExpireTime = 1200

// 未登录token超时时间
const UnLoginToKenExpireTime = UserLoginToKenExpireTime

// 响应体是否加密
const ReponseBodyCrypto = "reponseBodyCrypto"
