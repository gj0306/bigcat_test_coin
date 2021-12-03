package blc

// Version 公链版本信息默认为0
const Version = byte(0x00)

// CheckSum 两次sha256(公钥hash)后截取的字节数量
const CheckSum = 4