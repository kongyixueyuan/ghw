package BLC

import "bytes"

type TXOutput struct {
	Ghw_Value int64
	Ghw_Ripemd160Hash []byte  //用户名
}

func (txOutput *TXOutput)  Ghw_Lock(address string)  {

	publicKeyHash := Ghw_Base58Decode([]byte(address))

	txOutput.Ghw_Ripemd160Hash = publicKeyHash[1:len(publicKeyHash) - 4]
}


func Ghw_NewTXOutput(value int64,address string) *TXOutput {

	txOutput := &TXOutput{value,nil}

	// 设置Ripemd160Hash
	txOutput.Ghw_Lock(address)

	return txOutput
}

// 解锁
func (txOutput *TXOutput) Ghw_UnLockScriptPubKeyWithAddress(address string) bool {

	publicKeyHash := Ghw_Base58Decode([]byte(address))

	hash160 := publicKeyHash[1:len(publicKeyHash) - 4]

	return bytes.Compare(txOutput.Ghw_Ripemd160Hash,hash160) == 0
}



