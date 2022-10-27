// Package code 主要用于编解码string
// string通过 Encode 编码后变成 big.Int, 可以运用于之后的加密算法
// 解密结果是一个 big.Int, 可以用 Decode 解码成string，完整恢复原秘密/*
package code
