package main

import (
	"context"
	"ethereum-learn/util"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

// 教学视频：https://www.youtube.com/watch?v=nuivtRUaSw8&list=PLay9kDOVd_x7hbhssw4pTKZHzzc6OG0e_&index=4
var (
	arbitrumTestUrl = "https://arbitrum-sepolia.infura.io/v3/37ef3e9c7aff4830aae77fa3746ccb37"
	etherTestUrl    = "https://sepolia.infura.io/v3/37ef3e9c7aff4830aae77fa3746ccb37"
	etherMainUrl    = "https://mainnet.infura.io/v3/37ef3e9c7aff4830aae77fa3746ccb37"
)

func main() {
	client, err := ethclient.Dial(etherTestUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	a1 := common.HexToAddress("0x9f48b812a9aa300e195514805a3321d9fc870122")
	a2 := common.HexToAddress("0x039bf69e125d3abacd8b4404004fcf8d38b53c53")
	prvKey1, err := crypto.HexToECDSA("b753c223ab20fc9052878217ddeda18c8d1c87735e2ccfc374ca8e9c541c6318")
	if err != nil {
		log.Fatal(err)
	}

	b1, err := client.BalanceAt(context.Background(), a1, nil)
	if err != nil {
		log.Fatal(err)
	}

	b2, err := client.BalanceAt(context.Background(), a2, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("balance 1:", b1)
	fmt.Println("balance 2:", b2)

	nonce, err := client.PendingNonceAt(context.Background(), a1)
	if err != nil {
		log.Fatal(err)
	}

	//1 ether = 100000000000000000 wei
	//amount := big.NewInt(10000000000000000)
	amount := util.ToWei(0.01, 18) //发送0.01 eth

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(nonce, a2, amount, 21000, gasPrice, nil)

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), prvKey1)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent:%s\n", signTx.Hash().Hex())
}
