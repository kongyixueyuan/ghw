package main

import (
	"block/powdemo/blc"
	"fmt"
)

func main() {
	block:=blc.NewBlock("test hash",1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	fmt.Printf("%d\n",block.Nonce)
	fmt.Printf("%x\n",block.Hash)

	proofOfWork := blc.NewProofOfWork(block)

	fmt.Printf("%v",proofOfWork.IsValild())
}










