package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
)

//教学视频：https://www.youtube.com/watch?v=nuivtRUaSw8&list=PLay9kDOVd_x7hbhssw4pTKZHzzc6OG0e_&index=2

// infura 节点
var infuraURL = "https://mainnet.infura.io/v3/37ef3e9c7aff4830aae77fa3746ccb37"

// 本地测试节点
var ganacheURL = "http://127.0.0.1:8545"

func main() {
	//获取主网节点区块数量
	client, err := ethclient.DialContext(context.Background(), infuraURL)
	if err != nil {
		log.Fatalf("Error to create a ether client:%v", err)
	}
	defer client.Close()
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error to get a block:%v", err)
	}
	fmt.Printf("The block number:%v\n", block.Number())

	//获取钱包信息
	addr := "0x3eCD586c2eB666ECd32D6C54c5466aac30a874e8"
	address := common.HexToAddress(addr)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("Error to get a balance:%v", err)
	}
	fmt.Printf("The balance:%v\n", balance)
	//1 ether = 10^18 wei
	fBlance := new(big.Float)
	fBlance.SetString(balance.String())
	// 10*10*10*10*...18
	balanceEther := new(big.Float).Quo(fBlance, big.NewFloat(math.Pow10(18)))
	fmt.Printf("The ether balance:%v\n", balanceEther)
}
