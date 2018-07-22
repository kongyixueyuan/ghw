package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"strings"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
)

// 创世区块，Token数量
const subsidy  = 10

type Ghw_Transaction struct {
	Ghw_ID   []byte
	Ghw_Vin  []Ghw_TXInput
	Ghw_Vout []Ghw_TXOutput
}

// 是否是创世区块交易
func (tx Ghw_Transaction) Ghw_IsCoinbase() bool {
	// Vin 只有一条
	// Vin 第一条数据的Txid 为 0
	// Vin 第一条数据的Vout 为 -1
	return len(tx.Ghw_Vin) == 1 && len(tx.Ghw_Vin[0].Ghw_Txid) == 0 && tx.Ghw_Vin[0].Ghw_Vout == -1
}


// 将交易进行Hash
func (tx *Ghw_Transaction) Ghw_Hash() []byte  {
	var hash [32]byte

	txCopy := *tx
	txCopy.Ghw_ID = []byte{}

	hash = sha256.Sum256(txCopy.Ghw_Serialize())
	return hash[:]
}
// 新建创世区块的交易
func Ghw_NewCoinbaseTX(to ,data string) *Ghw_Transaction  {
	if data == ""{
		//如果数据为空，可以随机给默认数据,用于挖矿奖励
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}
	txin := Ghw_TXInput{[]byte{},-1,nil,[]byte(data)}
	txout := Ghw_NewTXOutput(subsidy,to)

	tx := Ghw_Transaction{nil,[]Ghw_TXInput{txin},[]Ghw_TXOutput{*txout}}
	tx.Ghw_ID = tx.Ghw_Hash()
	return &tx
}

// 转帐时生成交易
func Ghw_NewUTXOTransaction(wallet *Ghw_Wallet,to string,amount int,UTXOSet *Ghw_UTXOSet,txs []*Ghw_Transaction) *Ghw_Transaction   {

	// 如果本区块中，多笔转账
	/**
	第一种情况：
	  A:10
	  A->B 2
	  A->C 4

	  tx1:
	      Vin:
	           ATxID  out ...
	      Vout:
	           A : 8
	           B : 2
	  tx1:
	      Vin:
	           ATxID  out ...
	      Vout:
	           A : 4
	           C : 4
	第二种情况：
	  A:10+10
	  A->B 4
	  A->C 8
	**/

	pubKeyHash := Ghw_HashPubKey(wallet.Ghw_PublicKey)
	if len(txs) > 0 {
		// 查的txs中的UTXO
		utxo := Ghw_FindUTXOFromTransactions(txs)

		// 找出当前钱包已经花费的
		unspentOutputs := make(map[string][]int)
		acc := 0
		for txID,outs := range utxo {
			for outIdx, out := range outs.Ghw_Outputs {
				if out.Ghw_IsLockedWithKey(pubKeyHash) && acc < amount {
					acc += out.Ghw_Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		if acc >= amount { // 当前交易中的剩余余额可以支付
			fmt.Println("txs>0 && acc >= amount")
			return Ghw_NewUTXOTransactionEnd(wallet,to,amount,UTXOSet,acc,unspentOutputs,txs)
		}else{
			fmt.Println("txs>0 && acc < amount")
			accLeft, validOutputs := UTXOSet.Ghw_FindSpendableOutputs(pubKeyHash,  amount - acc)
			for k,v := range unspentOutputs{
				validOutputs[k] = v
			}
			return Ghw_NewUTXOTransactionEnd(wallet,to,amount,UTXOSet,acc + accLeft,validOutputs,txs)
		}
	} else { //只是当前一笔交易
		fmt.Println("txs==0")
		acc, validOutputs := UTXOSet.Ghw_FindSpendableOutputs(pubKeyHash, amount)

		return Ghw_NewUTXOTransactionEnd(wallet,to,amount,UTXOSet,acc,validOutputs,txs)
	}
}

func Ghw_NewUTXOTransactionEnd(wallet *Ghw_Wallet,to string,amount int,UTXOSet *Ghw_UTXOSet,acc int,UTXO map[string][]int,txs []*Ghw_Transaction) *Ghw_Transaction {

	if acc < amount {
		log.Panic("账户余额不足")
	}

	var inputs []Ghw_TXInput
	var outputs []Ghw_TXOutput
	// 构造input
	for txid, outs := range UTXO {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := Ghw_TXInput{txID, out, nil, wallet.Ghw_PublicKey}
			inputs = append(inputs, input)
		}
	}
	// 生成交易输出
	outputs = append(outputs, *Ghw_NewTXOutput(amount, to))
	// 生成余额
	if acc > amount {
		outputs = append(outputs, *Ghw_NewTXOutput(acc-amount, string(wallet.Ghw_GetAddress())))
	}

	tx := Ghw_Transaction{nil, inputs, outputs}
	tx.Ghw_ID = tx.Ghw_Hash()
	// 签名

	//tx.String()
	UTXOSet.Ghw_Blockchain.Ghw_SignTransaction(&tx, wallet.Ghw_PrivateKey,txs)

	return &tx
}


// 找出交易中的utxo
func Ghw_FindUTXOFromTransactions(txs []*Ghw_Transaction) map[string]Ghw_TXOutputs {
	UTXO := make(map[string]Ghw_TXOutputs)
	// 已经花费的交易txID : TXOutputs.index
	spentTXOs := make(map[string][]int)
	// 循环区块中的交易
	for _, tx := range txs {
		// 将区块中的交易hash，转为字符串
		txID := hex.EncodeToString(tx.Ghw_ID)

	Outputs:
		for outIdx, out := range tx.Ghw_Vout { // 循环交易中的 TXOutputs
			// Was the output spent?
			// 如果已经花费的交易输出中，有此输出，证明已经花费
			if spentTXOs[txID] != nil {
				for _, spentOutIdx := range spentTXOs[txID] {
					if spentOutIdx == outIdx { // 如果花费的正好是此笔输出
						continue Outputs // 继续下一次循环
					}
				}
			}

			outs := UTXO[txID] // 获取UTXO指定txID对应的TXOutputs
			outs.Ghw_Outputs = append(outs.Ghw_Outputs, out)
			UTXO[txID] = outs
		}

		if tx.Ghw_IsCoinbase() == false { // 非创世区块
			for _, in := range tx.Ghw_Vin {
				inTxID := hex.EncodeToString(in.Ghw_Txid)
				spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Ghw_Vout)
			}
		}
	}
	return UTXO

}

// 签名
func (tx *Ghw_Transaction) Ghw_Sign(privateKey ecdsa.PrivateKey,prevTXs map[string]Ghw_Transaction)  {
	if tx.Ghw_IsCoinbase() { // 创世区块不需要签名
		return
	}

	// 检查交易的输入是否正确
	for _,vin := range tx.Ghw_Vin{
		if prevTXs[hex.EncodeToString(vin.Ghw_Txid)].Ghw_ID == nil{
			log.Panic("错误：之前的交易不正确")
		}
	}

	txCopy := tx.Ghw_TrimmedCopy()

	for inID, vin := range txCopy.Ghw_Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Ghw_Txid)]
		txCopy.Ghw_Vin[inID].Ghw_Signature = nil
		txCopy.Ghw_Vin[inID].Ghw_PubKey = prevTx.Ghw_Vout[vin.Ghw_Vout].Ghw_PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Ghw_Vin[inID].Ghw_Signature = signature
		txCopy.Ghw_Vin[inID].Ghw_PubKey = nil
	}
}
// 验证签名
func (tx *Ghw_Transaction) Ghw_Verify(prevTXs map[string]Ghw_Transaction) bool {
	if tx.Ghw_IsCoinbase() {
		return true
	}

	for _, vin := range tx.Ghw_Vin {
		if prevTXs[hex.EncodeToString(vin.Ghw_Txid)].Ghw_ID == nil {
			log.Panic("错误：之前的交易不正确")
		}
	}

	txCopy := tx.Ghw_TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Ghw_Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Ghw_Txid)]
		txCopy.Ghw_Vin[inID].Ghw_Signature = nil
		txCopy.Ghw_Vin[inID].Ghw_PubKey = prevTx.Ghw_Vout[vin.Ghw_Vout].Ghw_PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Ghw_Signature)
		r.SetBytes(vin.Ghw_Signature[:(sigLen / 2)])
		s.SetBytes(vin.Ghw_Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.Ghw_PubKey)
		x.SetBytes(vin.Ghw_PubKey[:(keyLen / 2)])
		y.SetBytes(vin.Ghw_PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Ghw_Vin[inID].Ghw_PubKey = nil
	}

	return true
}

// 复制交易（输入的签名和公钥置为空）
func (tx *Ghw_Transaction) Ghw_TrimmedCopy() Ghw_Transaction {
	var inputs []Ghw_TXInput
	var outputs []Ghw_TXOutput

	for _, vin := range tx.Ghw_Vin {
		inputs = append(inputs, Ghw_TXInput{vin.Ghw_Txid, vin.Ghw_Vout, nil, nil})
	}

	for _, vout := range tx.Ghw_Vout {
		outputs = append(outputs, Ghw_TXOutput{vout.Ghw_Value, vout.Ghw_PubKeyHash})
	}

	txCopy := Ghw_Transaction{tx.Ghw_ID, inputs, outputs}

	return txCopy
}
// 打印交易内容
func (tx Ghw_Transaction) String()  {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction ID: %x", tx.Ghw_ID))

	for i, input := range tx.Ghw_Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Ghw_Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Ghw_Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Ghw_Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.Ghw_PubKey))
	}

	for i, output := range tx.Ghw_Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Ghw_Value))
		lines = append(lines, fmt.Sprintf("       PubKeyHash: %x", output.Ghw_PubKeyHash))
	}
	fmt.Println(strings.Join(lines, "\n"))
}


// 将交易序列化
func (tx Ghw_Transaction) Ghw_Serialize() []byte  {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil{
		log.Panic(err)
	}
	return encoded.Bytes()
}
// 反序列化交易
func Ghw_DeserializeTransaction(data []byte) Ghw_Transaction {
	var transaction Ghw_Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}

// 将交易数组序列化
func Ghw_SerializeTransactions(txs []*Ghw_Transaction) [][]byte  {

	var txsHash [][]byte
	for _,tx := range txs{
		txsHash = append(txsHash, tx.Ghw_Serialize())
	}
	return txsHash
}

// 反序列化交易数组
func Ghw_DeserializeTransactions(data [][]byte) []Ghw_Transaction {
	var txs []Ghw_Transaction
	for _,tx := range data {
		var transaction Ghw_Transaction
		decoder := gob.NewDecoder(bytes.NewReader(tx))
		err := decoder.Decode(&transaction)
		if err != nil {
			log.Panic(err)
		}

		txs = append(txs, transaction)
	}
	return txs
}
