package BLC

import (
	"fmt"
	"log"
)

func (cli *Ghw_CLI) ghw_startNode(nodeID, minerAddress string) {
	fmt.Printf("启动节点： %s\n", nodeID)
	if len(minerAddress) > 0 {
		fmt.Printf("当前节点为挖矿节点,挖矿地址为：%s\n",minerAddress)
		if Ghw_ValidateAddress(minerAddress) {
			fmt.Println("挖矿开始，挖矿地址为: ", minerAddress)
		} else {
			log.Panic("挖矿地址错误!")
		}
	}
	Ghw_StartServer(nodeID, minerAddress)
}
