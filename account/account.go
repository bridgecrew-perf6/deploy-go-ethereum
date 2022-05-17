package account

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"reflect"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// Account is an instance of the etherium account
type Account struct {
	ctx     context.Context
	client  *ethclient.Client
	address common.Address
}

// Balance is a current balance of the etherium account
type Balance struct {
	Wei  *big.Int
	Gwei *big.Int
	ETH  *big.Int
}

// NewAccount activates and validates existing eth account
func NewAccount(ctx context.Context, client *ethclient.Client, addressStr string) (*Account, error) {

	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(addressStr) {
		return nil, fmt.Errorf("invalid address format")
	}

	address := common.HexToAddress(addressStr)
	bytecode, err := client.CodeAt(ctx, address, nil) // nil is latest block
	if err != nil {
		return nil, fmt.Errorf("invalid address with err: %v", err)
	}
	isContract := len(bytecode) > 0
	if isContract {
		return nil, fmt.Errorf("address smart contract not allowed")
	}
	isZero := isZeroAddress(addressStr)
	if isZero {
		return nil, fmt.Errorf("zero address not allowed")
	}

	return &Account{ctx: ctx, client: client, address: address}, nil
}

// NewWallet creates new eth account TBD: manage keys
func NewWallet(ctx context.Context, client *ethclient.Client) (*Account, error) {

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("crypto.GenerateKey err: %v", err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyStr := hexutil.Encode(privateKeyBytes)[2:]
	fmt.Println(privateKeyStr)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		{
			return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyStr := hexutil.Encode(publicKeyBytes)[4:]
	fmt.Println(publicKeyStr)

	addressStr := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	address := common.HexToAddress(addressStr)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))

	return &Account{ctx: ctx, client: client, address: address}, nil
}

// Balance returns wei balance of the account
func (a Account) Balance() (*Balance, error) {
	balance, err := a.client.BalanceAt(a.ctx, a.address, nil)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	gwei := new(big.Int).Div(balance, new(big.Int).SetInt64(1000000000))
	eth := new(big.Int).Div(balance, new(big.Int).SetInt64(1000000000000000000))

	return &Balance{Wei: balance, Gwei: gwei, ETH: eth}, nil
}

// isZeroAddress validate if it's a 0 address
func isZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

// isValidAddress validate hex address
func isValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}
