package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"log"
)

//教学视频：https://www.youtube.com/watch?v=nuivtRUaSw8&list=PLay9kDOVd_x7hbhssw4pTKZHzzc6OG0e_&index=4

// 未安装iis，无法启动，仅记录
func main() {
	keystore := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "password"
	account, err := keystore.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Account created:%v\n", account)
}
