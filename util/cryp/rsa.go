package cryp

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// 可通过openssl产生
//openssl genrsa -out rsa_private_key.pem 1024
//var privateKey = []byte(`
//-----BEGIN RSA PRIVATE KEY-----
//MIIEowIBAAKCAQEAzLeI8JGMWFSjJaOzDFxloLwMsH6vwPZASttwVX0U0eiNOf3X
///Q08+BbTJB1W2Xy4YQv7/Y1slI3ZpcKu6FoEowW/Nf4RmlKgMflENdjD6BijQUDn
//9qYYBEJRkLSMSD+EiKUfq08z5x9i08NZg7E4u10mColtSuFp4RoAQKpFEA5qyrYd
//KmXoyW3JuTyszwO+aCMAS1dXKsKBF4HbC0GPbt8/7Es2E3zHJasxkL2iZ5mh3awJ
//wVnz/I0w+TGAu/XGzshGBAKDxHIRFS82YPw7/zwpk8NquZ55bzyNxCXIjvtLp2R7
//A/HmYwJ2u/xh9AtQpEj0sXW95Ouo/OO+1QElNwIDAQABAoIBABU46Z9W11/I1mju
//gX9EjNyO4hnh6EJuxNd9zDVwlBn2q71ZTWzUVH+7jgPubrR5M3wMDAGLCbiUw/1l
//I1C/FD/6NopYXmbLLgRAPQv//r8u3q3DFskBCvhWD7KapPhQbWLlC1VtDoplPI+L
//btoyIxl5XJo3CPd8SselNGV/wU0aogCjTaa40oauC0wiw/WoeCRA2Ok4GsQrYAN0
//IbV0/8jIBzIQONbENVDncm/u5vss/r7br0e/jdzcP0reMofAdLcQDn4nsU1JYD4i
//Ex087EmV2HTnn+uew8wtZlmsAZxKZh2FyTa0fCiVQXMw3Lf1GiZLXCs43ZysHJI3
//DNpNL9ECgYEA70Xuql67jvwhYndNtifLBdszrncTiNnT+B0T6ketSBKRCdBLGgnz
//hDm971dQIFjko8aQwsecg7u6t2zz8vHDXFUk1yeeOqhYMYbzpkfnR2k9TQ4IPJtl
//xgWoQ9S80jPFrG2F0rvFdULhP8eE7UOqiPdP30TLZyJreEPOoYiVx7sCgYEA2wct
//2kPFvCqSyoAtz4ZsbQMQIgxLfkXZM3zy1TKPCDd2E/ey1LV4Og/+6vrnQZ8GAhcD
//1brbmO6syepilB897bI0hNu03dWwT49Z4t6tDQ8YVOudm5HzFu2kK4ZEvayF8Ot/
//mZFvevUe893/VZwJlW7RpxMZFmtRpV5wHB3+6rUCgYAJ+1LfjKAqcN47q1p0lOhl
//UCWxy4nnFZ9AJIZmKaNS9GNUk3nuliewhnAkAfJ3xv2Sz3/OgGFJJZW+fS8YHXnW
//6j5lM2PocolrV4Pmle1SD1PdWQ6C6MCwKCBC5CcUZdCDRvZkOi0cnTOkY4BqHX6J
//xDdyyv3pSYhONhXyqy4EbQKBgQDUIUredv8etBkBeU1lDasblXjdkQzY2mt3u48w
//v0vaSGTbB+6ypqMvkOhyytiJLKxT/9hd+yDOKHM/B/u7u9ptyUemWWf95gVhuNP0
//r3fpCvKk5KH710oZrcVvxhXzohEDegJWSI4xBxCYXiz6zCpYCUGSUCPfG8eyoxlv
//kfmfdQKBgHSvwUPAYcNATSSaXbAJErLqknsRUdq9EWTNqJ8W4WrkTx2dWjW6deE9
//7VfIdxj/zUBe5fMD9BQXejP+nZDOdJQyXjfDcNmCGR5AEKiG9CzRey/w+rH+CbKf
//0R6Scku/cP3itx9akIgV7NIyIxq0HQWn7opkfGH3rGR4TtGmCU7w
//-----END RSA PRIVATE KEY-----
//`)

//openssl
//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
//var publicKey = []byte(`
//-----BEGIN PUBLIC KEY-----
//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzLeI8JGMWFSjJaOzDFxl
//oLwMsH6vwPZASttwVX0U0eiNOf3X/Q08+BbTJB1W2Xy4YQv7/Y1slI3ZpcKu6FoE
//owW/Nf4RmlKgMflENdjD6BijQUDn9qYYBEJRkLSMSD+EiKUfq08z5x9i08NZg7E4
//u10mColtSuFp4RoAQKpFEA5qyrYdKmXoyW3JuTyszwO+aCMAS1dXKsKBF4HbC0GP
//bt8/7Es2E3zHJasxkL2iZ5mh3awJwVnz/I0w+TGAu/XGzshGBAKDxHIRFS82YPw7
///zwpk8NquZ55bzyNxCXIjvtLp2R7A/HmYwJ2u/xh9AtQpEj0sXW95Ouo/OO+1QEl
//NwIDAQAB
//-----END PUBLIC KEY-----
//`)

// 加密
func RsaEncrypt(origData []byte, publicKey string) ([]byte, error) {

	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte, privateKey string) ([]byte, error) {
	//解密
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// 公钥解密
func PubKeyDECRYPT(input []byte, publicKey string) ([]byte, error) {
	// 类型断言
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	output := bytes.NewBuffer(nil)
	err = pubKeyIO(pub, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// 私钥加密
func PriKeyENCTYPT(input []byte, privateKey string) ([]byte, error) {

	//解密
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	output := bytes.NewBuffer(nil)
	err = priKeyIO(priv, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// 生成rsa的密钥对, 返回字节切片
func GenerateRsaKey(keySize int) (privateKeyStr, publicKeyStr []byte) {
	// 1. 使用rsa中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		panic(err)
	}
	// 2. 通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	derText := x509.MarshalPKCS1PrivateKey(privateKey)
	// 3. 要组织一个pem.Block(base64编码)
	block := pem.Block{
		Type:  "Elliptic curve cryptography private 5120", // rsa private key这个地方写个字符串就行
		Bytes: derText,
	}
	// 4. pem编码
	// 定义一个字节缓冲, 快速地连接字符串
	buffer := new(bytes.Buffer)
	pem.Encode(buffer, &block)
	// 将连接好的字节数组转换为字符串并输出
	privateKeyStr = buffer.Bytes()
	// ============ 公钥 ==========
	// 1. 从私钥中取出公钥
	publicKey := privateKey.PublicKey
	// 2. 使用x509标准序列化
	derstream, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	// 3. 将得到的数据放到pem.Block中
	block = pem.Block{
		Type:  "Elliptic curve cryptography private 5120",
		Bytes: derstream,
	}

	buffer = new(bytes.Buffer)
	pem.Encode(buffer, &block)
	// 将连接好的字节数组转换为字符串并输出
	publicKeyStr = buffer.Bytes()
	return
}
