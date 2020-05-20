package cryp

import (
	"fmt"
	"testing"
)

func TestAesDecrypt(t *testing.T) {
	orig := "hello world"
	//key := "123456781234567812345678"
	key := "9871267812345mn812345xyz"
	fmt.Println("原文：", orig)
	encryptCode, _ := AesEncrypt(orig, key)
	fmt.Println("密文：", encryptCode)
	decryptCode, _ := AesDecrypt(encryptCode, key+"2345t")
	fmt.Println("解密结果：", decryptCode)
}
