package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Ghw_TXOutputs struct {
	Ghw_Outputs []Ghw_TXOutput
}

//  序列化 TXOutputs
func (outs Ghw_TXOutputs) Ghw_Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// 反序列化 TXOutputs
func Ghw_DeserializeOutputs(data []byte) Ghw_TXOutputs {
	var outputs Ghw_TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
