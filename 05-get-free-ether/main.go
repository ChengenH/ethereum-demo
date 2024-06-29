package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

// 教学视频：https://www.youtube.com/watch?v=nuivtRUaSw8&list=PLay9kDOVd_x7hbhssw4pTKZHzzc6OG0e_&index=4
var (
	arbitrumTestUrl = "https://arbitrum-sepolia.infura.io/v3/37ef3e9c7aff4830aae77fa3746ccb37"
	etherTestUrl    = "https://arbitrum-sepolia.infura.io/v3/37ef3e9c7aff4830aae77fa3746ccb37"
)

func main() {
	client, err := ethclient.Dial(etherTestUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	address1 := common.HexToAddress("0x3ecd586c2eb666ecd32d6c54c5466aac30a874e8")
	address2 := common.HexToAddress("0x039bf69e125d3abacd8b4404004fcf8d38b53c53")
	balance1, err := client.BalanceAt(context.Background(), address1, nil)
	if err != nil {
		log.Fatal(err)
	}
	balance2, err := client.BalanceAt(context.Background(), address2, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance 1:%v\n", balance1)
	fmt.Printf("Balance 2:%v\n", balance2)
}
