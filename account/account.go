package account

import (
	"context"
	"fmt"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Account is an instance of the etherium account
type Account struct {
	ctx     context.Context
	client  *ethclient.Client
	address common.Address
}

func NewAccount(ctx context.Context, client *ethclient.Client, rawAddress string) (*Account, error) {

	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(rawAddress) {
		return nil, fmt.Errorf("invalid address format")
	}

	address := common.HexToAddress(rawAddress)
	bytecode, err := client.CodeAt(ctx, address, nil) // nil is latest block
	if err != nil {
		return nil, fmt.Errorf("invalid address with err: %v", err)
	}
	isContract := len(bytecode) > 0
	if isContract {
		return nil, fmt.Errorf("address smart contract not allowed")
	}
	return &Account{ctx: ctx, client: client, address: address}, nil
}

func NewAccount(ctx context.Context, client *ethclient.Client, rawAddress string) (*Account, error) {

	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(rawAddress) {
		return nil, fmt.Errorf("invalid address format")
	}

	address := common.HexToAddress(rawAddress)
	bytecode, err := client.CodeAt(ctx, address, nil) // nil is latest block
	if err != nil {
		return nil, fmt.Errorf("invalid address with err: %v", err)
	}
	isContract := len(bytecode) > 0
	if isContract {
		return nil, fmt.Errorf("address smart contract not allowed")
	}
	return &Account{ctx: ctx, client: client, address: address}, nil
}

func (a Account) NewWallet() (string, error) {
	var newPrivateKey string

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return newPrivateKey, fmt.Errorf("crypto.GenerateKey err: %v", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	newPrivateKey = hexutil.Encode(privateKeyBytes)
	return newPrivateKey, nil
}

func (w Wallet) GetWeiBalance(address common.Address) (*big.Int, error) {
	balance, err := w.Client.BalanceAt(w.Ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("clinet BalanceAt err: %v", err)
	}
	return balance, nil
}

func (a Wallet) Transaction(address common.Address) (*big.Int, error) {
	balance, err := a.Client.BalanceAt(a.Ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("clinet BalanceAt, err: %v", err)
	}
	return balance, nil
}

func (a Wallet) SendEthTo(publicKey string) error {

	return nil

}
