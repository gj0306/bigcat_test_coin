package blc

import (
	"bytes"
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/cloudflare/cfssl/scan/crypto/sha512"
	"math/big"
	"bigcat_test_coin/tools"
	"strings"
)

type Address struct {
	//私钥
	PrivateKey *ecdsa.PrivateKey
	//公钥
	PublicKey  ecdsa.PublicKey
}


func NewRandAddress() *Address {
	//随机key，用于创建公钥和私钥
	randKey := make([]byte, ketLength)
	_, _ = rand.Read(randKey)
	addr,_ := newAddressByBytes(randKey)
	return addr
}
func newAddressByBytes(bys []byte) (address *Address,err error) {
	//创建公钥和私钥
	prk, puk, err := getEcdsaKey(bys)
	if err != nil {
		return nil,fmt.Errorf("地址错误 %s",err.Error())
	}
	return &Address{
		PrivateKey: prk,
		PublicKey:  puk,
	},nil
}
func LoadAddress(privateKey string)(*Address, error){
	if privateKey == ""{
		return nil,fmt.Errorf("privateKey 为空")
	}
	k := big.NewInt(0).SetBytes(tools.Decodeb58(privateKey))
	priv := new(ecdsa.PrivateKey)
	priv.Curve = getCurveByKetLength(ketLength)
	priv.D = k
	priv.Public()
	priv.PublicKey.X, priv.PublicKey.Y = priv.Curve.ScalarBaseMult(k.Bytes())
	return &Address{
		PrivateKey: priv,
		PublicKey:  priv.PublicKey,
	},nil
}

// GetPrivateKey 获取私钥地址
func (a *Address) GetPrivateKey()string{
	bys := a.PrivateKey.D.Bytes()
	pri := tools.Encodeb58(bys)
	if len(pri) == privateKeyLength-1{
		pri = tools.Encodeb58([]byte{0}) + pri
	}
	return pri
}

func (a *Address) GetPublicKey()(puk []byte){
	x := a.PublicKey.X.Bytes()
	y := a.PublicKey.Y.Bytes()
	xl := len(x)
	yl := len(y)
	for xl != yl{
		if xl < yl{
			x = append([]byte{0}, x...)
		}else {
			y = append([]byte{0}, y...)
		}
		xl = len(x)
		yl = len(y)
	}
	return append(x, y...)
}
// GetAddress 通过公钥获得地址
func (a *Address) GetAddress() string {
	//1.ripemd160(sha256(publickey))
	ripPubKey := GeneratePublicKeyHash(a.GetPublicKey())
	//2.最前面添加一个字节的版本信息获得 versionPublickeyHash
	versionPublickeyHash := append([]byte{Version}, ripPubKey[:]...)
	//3.sha256(sha256(versionPublickeyHash))  取最后四个字节的值
	tailHash := checkSumHash(versionPublickeyHash)
	//4.拼接最终hash versionPublickeyHash + checksumHash
	finalHash := append(versionPublickeyHash, tailHash...)
	//进行base58加密
	address := tools.Encodeb58(finalHash)
	return address
}

// serliazle 序列化
func (a *Address) serliazle() []byte {
	var result bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(a)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

// deserialize 反序列化
func (a *Address) deserialize(d []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	gob.Register(elliptic.P256())
	err := decoder.Decode(a)
	if err != nil {
		panic(err.Error())
		//logs.Logger.Panic(err.Error())
	}
}

//修改 ketLength 长度，也要需要改 privateKeyLength 长度
const ketLength = 42
const privateKeyLength = 44
func getCurveByKetLength(ketLength int)(curve elliptic.Curve){
	if ketLength > 521/8+8 {
		curve = elliptic.P521()
	} else if ketLength > 384/8+8 {
		curve = elliptic.P384()
	} else if ketLength > 256/8+8 {
		curve = elliptic.P256()
	} else if ketLength > 224/8+8 {
		curve = elliptic.P224()
	}
	return
}
func CreatePublicKey(puk []byte)ecdsa.PublicKey{
	index := len(puk)/2
	return ecdsa.PublicKey{
		Curve: getCurveByKetLength(ketLength),
		X:     big.NewInt(0).SetBytes(puk[:index]),
		Y:     big.NewInt(0).SetBytes(puk[index:]),
	}
}

// 通过一个随机key创建公钥和私钥  随机key至少为36位
func getEcdsaKey(randBys []byte) (*ecdsa.PrivateKey, ecdsa.PublicKey, error) {
	var err error
	var prk *ecdsa.PrivateKey
	var puk ecdsa.PublicKey
	var curve elliptic.Curve
	length := len(randBys)
	if length < 224/8 {
		err = errors.New("私钥长度太短，至少为36位！")
		return prk, puk, err
	}
	curve = getCurveByKetLength(length)
	prk, err = ecdsa.GenerateKey(curve,strings.NewReader(string(randBys)))
	if err != nil {
		return prk, puk, err
	}
	puk = prk.PublicKey
	return prk, puk, err

}


// 对text加密，text必须是一个hash值，例如md5、sha1等
// 使用私钥prk
// 使用随机熵增强加密安全，安全依赖于此熵，randsign
// 返回加密结果，结果为数字证书r、s的序列化后拼接，然后用hex转换为string
func sign(text []byte, prk *ecdsa.PrivateKey) ([]byte, error) {
	//randSign 随机熵 用于加密安全 至少36位
	randBys := make([]byte,36)
	_, _ = rand.Read(randBys)
	r, s, err := ecdsa.Sign(strings.NewReader(string(randBys)), prk, text)
	if err != nil {
		return nil, err
	}
	rt, err := r.MarshalText()
	if err != nil {
		return nil, err
	}
	st, err := s.MarshalText()
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	_, err = w.Write([]byte(string(rt) + "+" + string(st)))
	if err != nil {
		return nil, err
	}
	_ = w.Flush()
	return b.Bytes(), nil
}


// 证书分解 通过hex解码，分割成数字证书r，s
func getSign(signature []byte) (rint, sint big.Int, err error) {
	//byterun, err := hex.DecodeString(signature)
	//if err != nil {
	//	err = errors.New("decrypt error, " + err.Error())
	//	return
	//}
	r, err := gzip.NewReader(bytes.NewBuffer(signature))
	if err != nil {
		err = errors.New("decode error," + err.Error())
		return
	}
	defer r.Close()
	buf := make([]byte, 1024)
	count, err := r.Read(buf)
	if err != nil {
		fmt.Println("decode = ", err)
		err = errors.New("decode read error," + err.Error())
		return
	}
	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return

}

// 校验文本内容是否与签名一致
// 使用公钥校验签名和文本内容
func verify(text []byte, signature []byte, key ecdsa.PublicKey) (bool, error) {
	rint, sint, err := getSign(signature)
	if err != nil {
		return false, err
	}
	result := ecdsa.Verify(&key, text, &rint, &sint)
	return result, nil
}


func GeneratePublicKeyHash(publicKey []byte) []byte {
	sha384PubKey := sha512.Sum384(publicKey)
	r := tools.NewRipemd160()
	r.Reset()
	r.Write(sha384PubKey[:])
	ripPubKey := r.Sum(nil)
	return ripPubKey
}

func checkSumHash(versionPublicKeyHash []byte) []byte {
	versionPublicKeyHashSha1 := sha256.Sum256(versionPublicKeyHash)
	versionPublicKeyHashSha2 := sha256.Sum256(versionPublicKeyHashSha1[:])
	tailHash := versionPublicKeyHashSha2[:CheckSum]
	return tailHash
}

// IsVerifyAddress 判断是否是有效的地址
func IsVerifyAddress(address string) bool {
	fullHash := tools.Decodeb58(address)
	if len(fullHash) != 25 {
		return false
	}
	prefixHash := fullHash[:len(fullHash)-CheckSum]
	tailHash := fullHash[len(fullHash)-CheckSum:]
	tailHash2 := checkSumHash(prefixHash)
	if bytes.Compare(tailHash, tailHash2[:]) == 0 {
		return true
	} else {
		return false
	}
}

// GetAddressFromPublicKey 通过公钥信息获得地址
func GetAddressFromPublicKey(publicKey []byte) string {
	if publicKey == nil {
		return ""
	}
	b := Address{PublicKey: CreatePublicKey(publicKey)}
	return b.GetAddress()
}

// EllipticCurveSign 使用私钥进行数字签名
func EllipticCurveSign(privateKey *ecdsa.PrivateKey, hash []byte) []byte {
	signature,err := sign(hash,privateKey)
	if err != nil{
		fmt.Println("数字前面失败： ",err.Error())
	}
	return signature
}

// EllipticCurveVerify 使用公钥进行签名验证
func EllipticCurveVerify(publicKey []byte, signature []byte, hash []byte) bool {
	ok,err := verify(hash,signature, CreatePublicKey(publicKey))
	if err != nil{
		fmt.Println("签名验证失败： ",err.Error())
	}
	return ok
}