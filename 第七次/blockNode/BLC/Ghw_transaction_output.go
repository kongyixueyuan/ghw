package BLC

import "bytes"

type Ghw_TXOutput struct {
	Ghw_Value  int
	Ghw_PubKeyHash []byte
}
// 根据地址获取 PubKeyHash
func (out *Ghw_TXOutput) Ghw_Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.Ghw_PubKeyHash = pubKeyHash
}

// 判断是否当前公钥对应的交易输出(是否是某个人的交易输出)
func (out *Ghw_TXOutput) Ghw_IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.Ghw_PubKeyHash, pubKeyHash) == 0
}

func Ghw_NewTXOutput(value int, address string) *Ghw_TXOutput {
	txo := &Ghw_TXOutput{value, nil}
	txo.Ghw_Lock([]byte(address))
	return txo
}


