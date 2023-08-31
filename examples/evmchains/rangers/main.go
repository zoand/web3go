package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/joho/godotenv/autoload"
	web3 "github.com/zoand/web3go"
)

var (
	devNetwork = "https://robin.rangersprotocol.com/api/jsonrpc"
)

func main() {

	// change to your rpc provider
	var rpcProvider = devNetwork
	//auto set chain id
	web3, err := web3.NewWeb3(rpcProvider)
	if err != nil {
		panic(err)
	}
	fmt.Println("current chaingID: ", web3.Eth.GetChainId())
	blockNumber, err := web3.Eth.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current block number: ", blockNumber)

	// only for test
	privateKey := os.Getenv("eth_privateKey") // hex string format
	if len(privateKey) == 0 {
		panic("please replace to your privateKey and keep safe by yourself")
	}
	// setup your privateKey
	if err := web3.Eth.SetAccount(privateKey); err != nil {
		panic(err)
	}
	privateKeyData, err := hex.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}
	ecdsaPrivateKey, err := crypto.ToECDSA(privateKeyData)
	if err != nil {
		panic(err)
	}

	addr := crypto.PubkeyToAddress(ecdsaPrivateKey.PublicKey)
	fmt.Printf("Address %s\n", addr)

	maticBalance, err := web3.Eth.GetBalance(addr, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("MATIC balance %v\n", maticBalance)
	nonce, err := web3.Eth.GetNonce(web3.Eth.Address(), nil)
	if err != nil {
		panic(err)
	}
	tx, err := web3.Eth.SyncSendEIP1559RawTransaction(
		addr,
		web3.Utils.ToWei("0.1"),
		nonce,
		21000,
		web3.Utils.ToGWei(25),
		web3.Utils.ToGWei(325),
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Send 0.1 MATIC to self tx %s\n", tx.TxHash)
}
