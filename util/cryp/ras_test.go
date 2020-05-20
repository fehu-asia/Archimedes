package cryp

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestRsaDecrypt(t *testing.T) {
	data, err := RsaEncrypt([]byte("hello world"), "123")
	fmt.Println(err, base64.StdEncoding.EncodeToString(data))
	origData, err := RsaDecrypt(data, "123")
	fmt.Println(err, string(origData))
}

func TestGenerateRsaKey(t *testing.T) {
	privateKeyByte, publicKeyByte := GenerateRsaKey(2048)

	privateKeyStr := string(privateKeyByte)
	fmt.Println("privateKeyStr:\n", privateKeyStr)
	fmt.Println()
	publicKeyStr := string(publicKeyByte)
	fmt.Println("publicKeyStr:\n", publicKeyStr)

	data, err := RsaEncrypt([]byte("hello world"), publicKeyStr)
	fmt.Println(err, base64.StdEncoding.EncodeToString(data))
	origData, err := RsaDecrypt(data, privateKeyStr)
	fmt.Println(err, string(origData))

}

func TestPriKeyENCTYPT(t *testing.T) {
	privateKeyByte, publicKeyByte := GenerateRsaKey(2048)
	privateKeyStr := string(privateKeyByte)
	publicKeyStr := string(publicKeyByte)
	data, err := PriKeyENCTYPT([]byte("hello world"), privateKeyStr)
	fmt.Println(err, base64.StdEncoding.EncodeToString(data))
	origData, err := PubKeyDECRYPT(data, publicKeyStr)
	fmt.Println(err, string(origData))
}
