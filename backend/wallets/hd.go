package wallets

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// DeriveAddress 从助记词和索引生成子地址
func DeriveAddress(mnemonic string, index uint32) (string, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", index))
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", err
	}
	return account.Address.Hex(), nil
}
func DerivePrivateKey(mnemonic string, index uint32) (string, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", index))
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", err
	}
	privKey, err := wallet.PrivateKey(account)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%x", crypto.FromECDSA(privKey)), nil
}
