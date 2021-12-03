package blc

import (
	"crypto/sha256"
	"fmt"
	"bigcat_test_coin/tools"
	"testing"
	"time"
)

func TestGetAddress(t *testing.T) {
	t.Log("测试获取coin币公私钥，并生成地址")
	{
		for i:=0;i<1;i++{
			bicAddr := NewRandAddress()
			privateKey := bicAddr.GetPrivateKey()
			address := bicAddr.GetAddress()
			t.Log("\t私钥为：", privateKey)
			t.Log("\t地址为：", address)
			t.Log("\t地址格式是否正确：", IsVerifyAddress(address))
		}

	}

	{
		addr,_ := LoadAddress("9d7uj6PkeHMxwiLv5fzzZQWRuSuXXToLyjf3F1NGdrtv")
		privateKey := addr.GetPrivateKey()
		address := addr.GetAddress()
		t.Log("\t私钥为：", privateKey)
		t.Log("\t地址为：", address)
		t.Log("\t地址格式是否正确：", IsVerifyAddress(address))
	}

}

func TestSign(t *testing.T) {
	t.Log("测试数字签名是否可用")
	{
		bicAddr := NewRandAddress()
		privateKey := bicAddr.GetPrivateKey()
		address := bicAddr.GetAddress()
		t.Log("\t私钥为：", privateKey)
		t.Log("\t地址为：", address)
		//
		hash := sha256.Sum256(tools.Int64ToBytes(time.Now().UnixNano()))
		fmt.Printf("\t签名hash:%x\n签名hash长度:%d\n", hash, len(hash))
		signature := EllipticCurveSign(bicAddr.PrivateKey, hash[:])
		verifyhash := append(hash[:], []byte("\t知道为什么这么长的验证信息也会通过吗？因为这个椭圆曲线只验证信息的前256位也就是前32字节！！！根据当时传入的elliptic.P256()有关！！！！")...)
		fmt.Printf("\t验证hash:%x\n验证hash长度:%d\n:", verifyhash, len(verifyhash))
		if EllipticCurveVerify(bicAddr.GetPublicKey(), signature, verifyhash) {
			t.Log("\t签名信息验证通过")
		} else {
			t.Fatal("\t签名信息验证失败！！！")
		}
	}
}