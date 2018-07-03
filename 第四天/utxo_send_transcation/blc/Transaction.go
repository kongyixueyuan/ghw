package blc

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

//utxo ,交易模型
type Transaction struct {
	//1,交易hash
	TxHash []byte
	//2,交易输入
	Vins []*TXInput
	//3,交易输出
	Vouts []*TXOutput
}
//1,创建创世区块创建时的Transaction
func NewConinbaseTransaction(address string) *Transaction  {
	//消费记录
	tXInput := &TXInput{[]byte{},-1,"创世区块"}
	//收入
	tXOutput := &TXOutput{10,"gaohongwei"}
	txCoinbase :=&Transaction{[]byte{},[]*TXInput{tXInput},[]*TXOutput{tXOutput}}
	//设置hash值
	txCoinbase.HashTransaction()
	return txCoinbase
}

//判断是否是创世区块的transaction，传世的tx
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

//转账的时候的transaction
func NewSimpleTransaction(from string,to string,amount int) *Transaction {
	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	//代表消费
	bytes ,_ := hex.DecodeString("1b5032e0cf4851f84dd89b9154912c082e28d5aa7f141625a0641c8a74f61802")
	txInput := &TXInput{bytes,0,from}
	//fmt.Printf("s:%s\n",s)
	// 消费
	txIntputs = append(txIntputs,txInput)
	// 转账
	txOutput := &TXOutput{int64(amount),to}
	txOutputs = append(txOutputs,txOutput)
	// 找零
	txOutput = &TXOutput{4 - int64(amount),from}
	txOutputs = append(txOutputs,txOutput)
	tx := &Transaction{[]byte{},txIntputs,txOutputs}
	//设置hash值
	tx.HashTransaction()



	return tx
}

func (tx *Transaction) HashTransaction()  {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil{
		log.Panic(err)
	}
	hash :=sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

