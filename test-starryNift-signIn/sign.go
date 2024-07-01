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
	BNBMainUrl      = "https://bsc-mainnet.nodereal.io/v1/6370e32a6a574d91b7421e31f464319a"
)

func main() {
	client, err := ethclient.Dial(BNBMainUrl)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xE3bA0072d1da98269133852fba1795419D72BaF4") //签到合约地址

	// 加载合约ABI
	tokenSignInABI, err := abi.JSON(strings.NewReader(`[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"user","type":"address"},{"indexed":false,"internalType":"uint256","name":"date","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"consecutiveDays","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"rewardXP","type":"uint256"}],"name":"UserSignedIn","type":"event"},{"inputs":[],"name":"admin","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getConsecutiveSignInDays","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getLastSignInDate","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"day","type":"uint256"}],"name":"getRewardXP","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getSignInInterval","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getTimeUntilNextSignIn","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"hasBrokenStreak","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"setDefaultRewards","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_maxDays","type":"uint256"}],"name":"setMaxConsecutiveDays","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_interval","type":"uint256"}],"name":"setSignInInterval","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"day","type":"uint256"},{"internalType":"uint256","name":"rewardXP","type":"uint256"}],"name":"setSignInReward","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"signIn","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"day","type":"uint256"},{"internalType":"uint256","name":"reward","type":"uint256"}],"name":"updateSignReward","outputs":[],"stateMutability":"nonpayable","type":"function"}]`))
	if err != nil {
		log.Fatalf("Failed to parse token ABI: %v", err)
	}

	// 创建签到调用数据
	data, err := tokenSignInABI.Pack("signIn")
	if err != nil {
		log.Fatalf("Failed to pack transfer  %v", err)
	}

	//私钥
	privateKey, err := crypto.HexToECDSA("*********")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 签到操作
	err = signIn(client, fromAddress, contractAddress, data, privateKey)
	if err != nil {
		log.Printf("Failed to sign in: %v", err)
	} else {
		log.Println("Successfully signed in!")
	}

}

func signIn(client *ethclient.Client, fromAddress common.Address, contractAddress common.Address, data []byte, privateKey *ecdsa.PrivateKey) error {
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

	// 设置调用消息
	msg := ethereum.CallMsg{
		From: fromAddress,
		To:   &contractAddress,
		Data: data,
	}

	// 估算GasLimit
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatalf("Failed to estimate gas limit: %v", err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &contractAddress, // 注意这里是交互的合约地址
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
	fmt.Printf("tx is confirmed: %v. SignIn sucess. Block number: %v\n", receipt.Status, receipt.BlockNumber)

	return nil
}
