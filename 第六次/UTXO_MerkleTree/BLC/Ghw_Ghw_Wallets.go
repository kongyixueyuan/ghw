package BLC

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"io/ioutil"
	"log"
	"os"
)

const ghw_walletFile  = "Wallets.dat"

type Wallets struct {
	Ghw_WalletsMap map[string]*Wallet
}



// 创建钱包集合
func Ghw_NewWallets() (*Wallets,error){

	if _, err := os.Stat(ghw_walletFile); os.IsNotExist(err) {

		wallets := &Wallets{}

		wallets.Ghw_WalletsMap = make(map[string]*Wallet)

		return wallets,err
	}

	fileContent, err := ioutil.ReadFile(ghw_walletFile)

	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets

	gob.Register(elliptic.P256())

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))

	err = decoder.Decode(&wallets)

	if err != nil {
		log.Panic(err)
	}

	return &wallets,nil
}

// 创建一个新钱包
func (w *Wallets) Ghw_CreateNewWallet()  {

	wallet := Ghw_NewWallet()

	fmt.Printf("Address：%s\n",wallet.Ghw_GetAddress())

	w.Ghw_WalletsMap[string(wallet.Ghw_GetAddress())] = wallet

	w.Ghw_SaveWallets()
}

// 将钱包信息写入到文件
func (w *Wallets) Ghw_SaveWallets()  {

	var content bytes.Buffer

	// 注册的目的，是为了，可以序列化任何类型
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)

	err := encoder.Encode(&w)

	if err != nil {
		log.Panic(err)
	}

	// 将序列化以后的数据写入到文件，原来文件的数据会被覆盖
	err = ioutil.WriteFile(ghw_walletFile, content.Bytes(), 0644)

	if err != nil {
		log.Panic(err)
	}


}

