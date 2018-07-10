package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"

	"math/big"
	"crypto/elliptic"
	"time"
)

// UTXO
type Transaction struct {

	//1. 交易hash
	Ghw_TxHash []byte

	//2. 输入
	Ghw_Vins []*TXInput

	//3. 输出
	Ghw_Vouts []*TXOutput
}

//[]byte{}

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) Ghw_IsCoinbaseTransaction() bool {

	return len(tx.Ghw_Vins[0].Ghw_TxHash) == 0 && tx.Ghw_Vins[0].Ghw_Vout == -1
}

//1. Transaction 创建分两种情况
//1. 创世区块创建时的Transaction
func Ghw_NewCoinbaseTransaction(address string) *Transaction {

	//代表消费
	txInput := &TXInput{[]byte{},-1,nil,[]byte{}}

	txOutput := Ghw_NewTXOutput(10,address)

	txCoinbase := &Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}

	//设置hash值
	txCoinbase.Ghw_HashTransaction()

	return txCoinbase
}

// 事务hash
func (tx *Transaction) Ghw_HashTransaction()  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	resultBytes := bytes.Join([][]byte{Ghw_IntToHex(time.Now().Unix()),result.Bytes()},[]byte{})

	hash := sha256.Sum256(resultBytes)

	tx.Ghw_TxHash = hash[:]
}

//2. 转账时产生的Transaction
func Ghw_NewSimpleTransaction(from string,to string,amount int64,utxoSet *UTXOSet,txs []*Transaction) *Transaction {

	//$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	wallets,_ := Ghw_NewWallets()
	wallet := wallets.Ghw_WalletsMap[from]

	// 通过一个函数，返回
	money,spendableUTXODic := utxoSet.Ghw_FindSpendableUTXOS(from,amount,txs)
	//
	//	{hash1:[0],hash2:[2,3]}

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {

			txInput := &TXInput{txHashBytes,index,nil,wallet.Ghw_PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := Ghw_NewTXOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = Ghw_NewTXOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.Ghw_HashTransaction()

	//进行签名
	utxoSet.Ghw_Blockchain.Ghw_SignTransaction(tx, wallet.Ghw_PrivateKey,txs)

	return tx

}

//产生Hash
func (tx *Transaction) Ghw_Hash() []byte {

	txCopy := tx

	txCopy.Ghw_TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Ghw_Serialize())

	return hash[:]
}

//序列化
func (tx *Transaction) Ghw_Serialize() []byte {

	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)

	err := enc.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// 签名
func (tx *Transaction) Ghw_Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//创世区块不用签名
	if tx.Ghw_IsCoinbaseTransaction() {
		return
	}

	for _, vin := range tx.Ghw_Vins {
		if prevTXs[hex.EncodeToString(vin.Ghw_TxHash)].Ghw_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.Ghw_TrimmedCopy()

	for inID, vin := range txCopy.Ghw_Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Ghw_TxHash)]
		txCopy.Ghw_Vins[inID].Ghw_Signature = nil
		txCopy.Ghw_Vins[inID].Ghw_PublicKey = prevTx.Ghw_Vouts[vin.Ghw_Vout].Ghw_Ripemd160Hash
		txCopy.Ghw_TxHash = txCopy.Ghw_Hash()
		txCopy.Ghw_Vins[inID].Ghw_PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.Ghw_TxHash)
		if err != nil {
			log.Panic(err)
		}

		signature := append(r.Bytes(), s.Bytes()...)

		tx.Ghw_Vins[inID].Ghw_Signature = signature
	}
}

// 拷贝一份新的Transaction用于签名
func (tx *Transaction) Ghw_TrimmedCopy() Transaction {

	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Ghw_Vins {
		inputs = append(inputs, &TXInput{vin.Ghw_TxHash, vin.Ghw_Vout, nil, nil})
	}

	for _, vout := range tx.Ghw_Vouts {
		outputs = append(outputs, &TXOutput{vout.Ghw_Value, vout.Ghw_Ripemd160Hash})
	}

	txCopy := Transaction{tx.Ghw_TxHash, inputs, outputs}

	return txCopy
}

// 数字签名验证
func (tx *Transaction) Ghw_Verify(prevTXs map[string]Transaction) bool {

	if tx.Ghw_IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Ghw_Vins {
		if prevTXs[hex.EncodeToString(vin.Ghw_TxHash)].Ghw_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.Ghw_TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Ghw_Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Ghw_TxHash)]
		txCopy.Ghw_Vins[inID].Ghw_Signature = nil
		txCopy.Ghw_Vins[inID].Ghw_PublicKey = prevTx.Ghw_Vouts[vin.Ghw_Vout].Ghw_Ripemd160Hash
		txCopy.Ghw_TxHash = txCopy.Ghw_Hash()
		txCopy.Ghw_Vins[inID].Ghw_PublicKey = nil

		// 私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Ghw_Signature)
		r.SetBytes(vin.Ghw_Signature[:(sigLen / 2)])
		s.SetBytes(vin.Ghw_Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.Ghw_PublicKey)
		x.SetBytes(vin.Ghw_PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.Ghw_PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.Ghw_TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}
