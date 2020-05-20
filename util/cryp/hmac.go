package cryp

import (
	"crypto/hmac"
	"crypto/sha256"
)

// 注意：计算hamc_sha256时，是否需要转成十六进制，取决于自己的需要，代码中注释掉该行代码：
//hex.EncodeToString(h.Sum(nil))
//func ComputeHmacSha256(message string, secret string) string {
//	key := []byte(secret)
//	h := hmac.New(sha256.New, key)
//	h.Write([]byte(message))
//	//	fmt.Println(h.Sum(nil))
//	sha := hex.EncodeToString(h.Sum(nil))
//	//	fmt.Println(sha)
//	//return hex.EncodeToString(h.Sum(nil))
//	return base64.StdEncoding.EncodeToString([]byte(sha))
//}

// 生成消息认证码
func GenerateHamc(plainText, key string) []byte {
	// 1.创建哈希接口, 需要指定使用的哈希算法, 和秘钥
	myhash := hmac.New(sha256.New, []byte(key))
	// 2. 给哈希对象添加数据
	myhash.Write([]byte(plainText))
	// 3. 计算散列值
	hashText := myhash.Sum(nil)
	return hashText
}

// 验证消息认证码
func VerifyHamc(plainText, key, hashText string) bool {
	// 1.创建哈希接口, 需要指定使用的哈希算法, 和秘钥
	myhash := hmac.New(sha256.New, []byte(key))
	// 2. 给哈希对象添加数据
	myhash.Write([]byte(plainText))
	// 3. 计算散列值
	hamc1 := myhash.Sum(nil)
	// 4. 两个散列值比较
	return hmac.Equal([]byte(hashText), hamc1)
}
