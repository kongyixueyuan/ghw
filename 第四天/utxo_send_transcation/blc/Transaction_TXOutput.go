package blc

type TXOutput struct {
	Money int64 //交易金额、
	ScriptPublicKey string //用户名
}

func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	return txOutput.ScriptPublicKey == address
}
