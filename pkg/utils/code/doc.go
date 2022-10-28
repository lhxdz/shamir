// Package code 主要用于编解码string
// string通过 EncodeSecret 编码后变成 big.Int, 可以运用于之后的加密算法
// 解密结果是一个 big.Int, 可以用 DecodeSecret 解码成string，完整恢复原秘密
// 加密后生成的密钥对(x, y)，可以通过 DecodeKey 生成x,y对应的string
// 同样的，x,y对应的string可以通过 EncodeKey 恢复成大整数，方便用于解密
// /*
package code
