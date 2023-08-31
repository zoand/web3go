package main

import (
	"fmt"

	web3 "github.com/zoand/web3go"
)

func main() {

	// change to your rpc provider
	var rpcProvider = "https://rpc.flashbots.net"
	web3, err := web3.NewWeb3(rpcProvider)
	if err != nil {
		panic(err)
	}
	blockNumber, err := web3.Eth.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current block number: ", blockNumber)
}
