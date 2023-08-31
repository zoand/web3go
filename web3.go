package web3

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zoand/web3go/eth"
	"github.com/zoand/web3go/rpc"
	"github.com/zoand/web3go/utils"
)

type Web3 struct {
	Eth   *eth.Eth
	Utils *utils.Utils
	c     *rpc.Client
	url   string
}

func NewWeb3(provider string) (*Web3, error) {
	return NewWeb3WithProxy(provider, "")
}

func NewWeb3WithProxy(provider, proxy string) (*Web3, error) {
	c, err := rpc.NewClient(provider, proxy)
	if err != nil {
		return nil, err
	}
	e := eth.NewEth(c)

	// providerLowerStr := strings.ToLower(provider)

	// if strings.Contains(providerLowerStr, "ropsten") {
	// 	e.SetChainId(3)
	// } else if strings.Contains(providerLowerStr, "kovan") {
	// 	e.SetChainId(42)
	// } else if strings.Contains(providerLowerStr, "rinkeby") {
	// 	e.SetChainId(4)
	// } else if strings.Contains(providerLowerStr, "goerli") {
	// 	e.SetChainId(5)
	// } else {
	// 	e.SetChainId(1)
	// }

	u := utils.NewUtils()
	w := &Web3{
		Eth:   e,
		Utils: u,
		c:     c,
		url:   provider,
	}
	chainID := w.GetChainId(w.url)
	fmt.Println("chainID:", chainID)
	e.SetChainId(chainID.Int64())

	// Default poll timeout 2 hours
	w.Eth.SetTxPollTimeout(7200)
	return w, nil
}

func (w *Web3) Version() (string, error) {
	var out string
	err := w.c.Call("web3_clientVersion", &out)
	return out, err
}

func (w *Web3) GetChainId(rpcurl string) *big.Int {
	client, err := ethclient.Dial(rpcurl)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return chainID
}
