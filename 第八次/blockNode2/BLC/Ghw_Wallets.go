package BLC

import (
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
	"fmt"
)

const walletFile  = "wallet_%s.dat"

type Ghw_Wallets struct {
	Ghw_Wallets map[string]*Ghw_Wallet
}

// 生成新的钱包
// 从数据库中读取，如果不存在
func Ghw_NewWallets(nodeID string)(*Ghw_Wallets,error)  {
	wallets := Ghw_Wallets{}
	wallets.Ghw_Wallets = make(map[string]*Ghw_Wallet)

	err := wallets.Ghw_LoadFromFile(nodeID)

	return &wallets,err
}
// 生成新的钱包地址列表
func (ws *Ghw_Wallets) Ghw_NewWallet() *Ghw_Wallet {
	wallet := Ghw_NewWallet()
	address := wallet.Ghw_GetAddress()
	ws.Ghw_Wallets[string(address)] = wallet
	return wallet
}
// 获取钱包地址
func (ws *Ghw_Wallets) Ghw_GetAddresses()[]string  {
	var addresses []string
	for address := range ws.Ghw_Wallets{
		addresses = append(addresses,address)
	}
	return addresses
}

// 根据地址获取钱包的详细信息
func (ws Ghw_Wallets) Ghw_GetWallet(address string) Ghw_Wallet {
	return *ws.Ghw_Wallets[address]
}

// 从数据库中读取钱包列表
func (ws *Ghw_Wallets)Ghw_LoadFromFile(nodeID string) error  {
	 walletFile := fmt.Sprintf(walletFile, nodeID)
	 if _,err := os.Stat(walletFile) ; os.IsNotExist(err){
	 	return err
	 }

	 fileContent ,err := ioutil.ReadFile(walletFile)
	 if err !=nil{
	 	log.Panic(err)
	 }

	 var wallets Ghw_Wallets
	 gob.Register(elliptic.P256())
	 decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	 err = decoder.Decode(&wallets)
	 if err !=nil{
	 	log.Panic(err)
	 }

	 ws.Ghw_Wallets = wallets.Ghw_Wallets

	 return nil
}

// 将钱包存到数据库中
func (ws *Ghw_Wallets)Ghw_SaveToFile(nodeID string)  {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err !=nil{
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile,content.Bytes(),0644)
	if err !=nil{
		log.Panic(err)
	}
}
// 打印所有钱包的余额
func (ws *Ghw_Wallets) Ghw_GetBalanceAll(nodeID string) map[string]int {
	addresses := ws.Ghw_GetAddresses()
	bc := Ghw_NewBlockchain(nodeID)
	defer bc.ghw_db.Close()
	UTXOSet := Ghw_UTXOSet{bc}

	result := make(map[string]int)
	for _,address := range addresses{
		if !Ghw_ValidateAddress(address) {
			result[address] = -1
		}
		balance := UTXOSet.Ghw_GetBalance(address)
		result[address] = balance
	}
	return result
}