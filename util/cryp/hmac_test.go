package cryp

import (
	"fmt"
	"testing"
)

func TestComputeHmacSha256(t *testing.T) {
	message := "hello world!"
	secret := "0933e54e76b24731a2d84b6b463ec04c"
	fmt.Println(GenerateHamc(message, secret))
}
func TestVerifyHamc(t *testing.T) {
	//src := []byte("在消息认证码中，需要发送者和接收者之间共享密钥，而这个密钥不能被主动攻击者Mallory获取。" +
	//	"如果这个密钥落入Mallory手中，则Mallory也可以计算出MAC值，从而就能够自由地进行篡改和伪装攻击，" +
	//	"这样一来消息认证码就无法发挥作用了。")
	//key := []byte("helloworld")

	message := "hello world!"
	secret := "0933e54e76b24731a2d84b6b463ec04c"
	hamc1 := GenerateHamc(message, secret)
	bl := VerifyHamc(message, secret, string(hamc1))
	//fmt.Printf("校验结果: %t\n", bl)
	fmt.Println(bl)
}
