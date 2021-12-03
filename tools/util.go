package tools

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
)



// Int64ToBytes int64转换成字节数组
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

// BytesToInt 字节数组转换为int
func BytesToInt(bys []byte) int {
	bytebuffer := bytes.NewBuffer(bys)
	var data int64
	_ = binary.Read(bytebuffer, binary.BigEndian, &data)
	return int(data)
}

// GenerateRealRandom 生成随机数
func GenerateRealRandom() int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000000000000000))
	if err != nil {
		fmt.Println(err)
	}
	return n.Int64()
}

// StringsGetMaxCount 字符串统计次数
func StringsGetMaxCount(strs []string) string{
	mp := make(map[string]int)
	for _,s := range strs{
		mp[s]+=1
	}
	var sign string
	for s,_ := range mp{
		if mp[s]>mp[sign]{
			sign = s
		}
	}
	return sign
}

func NumberQuoAndRem (buf1 []byte,divisor int)(numberQuo,numberRem int){
	number1 := big.NewInt(0).SetBytes(buf1)
	number2 := big.NewInt(int64(divisor))
	numberQuo = number1.Quo(number1,number2).BitLen()
	numberRem = number1.Rem(number1,number2).BitLen()
	return
}

func Md5(bys []byte) string {
	has := md5.Sum(bys)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}