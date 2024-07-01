package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

	privateKey, err := crypto.HexToECDSA("b753c223ab20fc9052878217ddeda18c8d1c87735e2ccfc374ca8e9c541c6318")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got gas price: %d\n", gasPrice)

	toAddress := common.HexToAddress("0x039bf69e125d3abacd8b4404004fcf8d38b53c53")
	tokenAddress := common.HexToAddress("0x3F4B6664338F23d2397c953f2AB4Ce8031663f80") //OKB TEST

	// 加载合约ABI
	tokenABI, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}]`))
	if err != nil {
		log.Fatalf("Failed to parse token ABI: %v", err)
	}

	// 准备调用合约的参数
	parsedAmount := new(big.Int)
	parsedAmount, ok = parsedAmount.SetString("100000000000000000", 10)
	if !ok {
		log.Fatalf("Invalid amount")
	}

	// 创建转账调用数据
	data, err := tokenABI.Pack("transfer", toAddress, parsedAmount)
	if err != nil {
		log.Fatalf("Failed to pack transfer  %v", err)
	}

	// 设置调用消息
	msg := ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	}

	// 估算GasLimit
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatalf("Failed to estimate gas limit: %v", err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &tokenAddress, // 注意这里是代币的合约地址
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex()) // tx sent: 0xd6362823b5d924f0dae344939977e02df48f454c590e1e9e7d297b403b7f0294

	// wait until transaction is confirmed
	var receipt *types.Receipt
	for {
		receipt, err = client.TransactionReceipt(context.Background(), signedTx.Hash())
		if err != nil {
			fmt.Println("tx is not confirmed yet")
			time.Sleep(5 * time.Second)
		}
		if receipt != nil {
			break
		}
	}
	// Status = 1 if transaction succeeded
	fmt.Printf("tx is confirmed: %v. Block number: %v\n", receipt.Status, receipt.BlockNumber)
}
