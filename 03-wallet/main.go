package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

//教学视频：https://www.youtube.com/watch?v=nuivtRUaSw8&list=PLay9kDOVd_x7hbhssw4pTKZHzzc6OG0e_&index=3

// 未安装iis，无法启动，仅记录
func main() {
	pvk, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Error to generateKey:%v", err)
	}

	//私钥
	pData := crypto.FromECDSA(pvk)
	fmt.Println(hexutil.Encode(pData))

	//公钥
	puData := crypto.FromECDSAPub(&pvk.PublicKey)
	fmt.Println(hexutil.Encode(puData))

	//钱包地址
	fmt.Println(crypto.PubkeyToAddress(pvk.PublicKey).Hex())
}
