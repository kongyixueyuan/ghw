package main

import (
	"block/boltDb/blc"
	"fmt"
	"log"
	"github.com/boltdb/bolt"
)

func main()  {
// 创建创世区块
	block:=blc.NewBlock("创世区块",1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	fmt.Println("原始区块")
	fmt.Printf("nonce值%d\n",block.Nonce)
	fmt.Printf("hash值 %x\n",block.Hash)

	bytes := block.Serialize()
	fmt.Printf("序列化之后的区块%d:",bytes)
	//反序列化
	block= blc.DeserializeBlock(bytes)
/*	fmt.Println("反序列化后的区块")
	fmt.Printf("反序列化之后的区块%d:",block)
	fmt.Printf("nonce值%d\n",block.Nonce)
	fmt.Printf("hash值%x\n",block.Hash)*/

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


}



