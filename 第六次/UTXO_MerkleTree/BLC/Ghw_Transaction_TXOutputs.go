package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutputs struct {
	Ghw_UTXOS []*UTXO
}

// 将区块序列化成字节数组
func (txOutputs *TXOutputs) Ghw_Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(txOutputs)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func Ghw_DeserializeTXOutputs(txOutputsBytes []byte) *TXOutputs {

	var txOutputs TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))

	err := decoder.Decode(&txOutputs)

	if err != nil {
		log.Panic(err)
	}

	return &txOutputs
}