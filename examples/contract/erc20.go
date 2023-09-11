package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/joho/godotenv/autoload"
	web3 "github.com/zoand/web3go"
	"github.com/zoand/web3go/apps/erc20"
	"github.com/zoand/web3go/types"
)

/*
{
    "BSC": "0xB82318f4cB5D04936A12e91148230064B19e03f8",
    "ETH": "0xFe4d81A83deF8ba75E5B9670C562D124ceDb3e94",
    "MAP": "0xFe4d81A83deF8ba75E5B9670C562D124ceDb3e94",
    "Rangers": "0x9c1CbFE5328DFB1733d59a7652D0A49228c7E12C"
}
*/

func test_approve_erc20() {
	rpcProvider := os.Getenv("eth_rpc_provider")
	admin_account := os.Getenv("eth_privateKey")
	usdt_address := os.Getenv("eth_usdt_rpg_address")
	bridge_address := os.Getenv("eth_bridge_rpg_address")

	web3, err := web3.NewWeb3(rpcProvider)
	if err != nil {
		panic(err)
	}
	web3.Eth.SetAccount(admin_account)
	usdt, err := erc20.NewERC20(web3, common.HexToAddress(usdt_address))
	if err != nil {
		panic(err)
	}
	usdt.SetConfirmation(1)
	bridge := common.HexToAddress(bridge_address)
	amount := web3.Utils.ToWei("1")
	gasprice := web3.Utils.ToGWei(50)
	tx, err := usdt.Approve(bridge, amount, gasprice, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("hash:", tx.Hex())
}

func main() {
	test_approve_erc20()
	return

	abiStr := `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`

	// change to your rpc provider
	var rpcProvider = "https://robin.rangersprotocol.com/api/jsonrpc"

	web3, err := web3.NewWeb3(rpcProvider)

	if err != nil {
		panic(err)
	}
	web3.Eth.SetAccount(os.Getenv("eth_privateKey"))
	// set default account by private key
	privateKey := os.Getenv("eth_privateKey")
	kovanChainId := int64(42)
	fmt.Printf("Chain id %v\n", kovanChainId)
	if err := web3.Eth.SetAccount(privateKey); err != nil {
		panic(err)
	}
	web3.Eth.SetChainId(kovanChainId)
	tokenAddr := "0x0F3A62dB02F743b549053cc8d538C65aB01E3618"
	contract, err := web3.Eth.NewContract(abiStr, tokenAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Contract address: ", contract.Address())

	totalSupply, err := contract.Call("totalSupply")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total supply %v\n", totalSupply)

	data, _ := contract.EncodeABI("balanceOf", web3.Eth.Address())
	fmt.Printf("%x\n", data)

	balance, err := contract.Call("balanceOf", web3.Eth.Address())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Balance of %v is %v\n", web3.Eth.Address(), balance)

	allowance, err := contract.Call("allowance", web3.Eth.Address(), "0x0000000000000000000000000000000000000002")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Allowance is %v\n", allowance)
	approveInputData, err := contract.Methods("approve").Inputs.Pack("0x0000000000000000000000000000000000000002", web3.Utils.ToWei("0.2"))
	if err != nil {
		panic(err)
	}
	// fmt.Println(approveInputData)

	tokenAddress := common.HexToAddress(tokenAddr)

	call := &types.CallMsg{
		From: web3.Eth.Address(),
		To:   tokenAddress,
		Data: approveInputData,
		Gas:  types.NewCallMsgBigInt(big.NewInt(types.MAX_GAS_LIMIT)),
	}
	// fmt.Printf("call %v\n", call)
	gasLimit, err := web3.Eth.EstimateGas(call)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Estimate gas limit %v\n", gasLimit)
	nonce, err := web3.Eth.GetNonce(web3.Eth.Address(), nil)
	if err != nil {
		panic(err)
	}
	txHash, err := web3.Eth.SyncSendRawTransaction(
		common.HexToAddress(tokenAddr),
		big.NewInt(0),
		nonce,
		gasLimit,
		web3.Utils.ToGWei(1),
		approveInputData,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Send approve tx hash %v\n", txHash)
}
