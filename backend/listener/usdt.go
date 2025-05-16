package listener

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/Ice-American/hdwallet-usdt-payment-listener/backend/models"
)

const USDTAddress = "0x55d398326f99059fF775485246999027B3197955"

func ListenUSDT() {
	client, err := ethclient.Dial("wss://bsc-rpc.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	transferSig := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

	// 取出所有 sub_address，组装过滤器
	rows, _ := models.DB.Query("SELECT sub_address FROM deposits WHERE TRUE")
	var subs []common.Address
	for rows.Next() {
		var sa string
		rows.Scan(&sa)
		subs = append(subs, common.HexToAddress(sa))
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(USDTAddress)},
		Topics:    [][]common.Hash{{transferSig}},
	}
	logsChan := make(chan types.Log)
	client.SubscribeFilterLogs(context.Background(), query, logsChan)

	for vLog := range logsChan {
		common.HexToAddress(vLog.Topics[1].Hex())
		to := common.HexToAddress(vLog.Topics[2].Hex())
		amt := new(big.Int).SetBytes(vLog.Data)

		// 检查是否是我们关注的子地址
		for _, sa := range subs {
			if to == sa {
				// 插入 payments
				_, err := models.DB.Exec(
					`INSERT INTO payments(deposit_id, tx_hash, amount)
           SELECT id, $1, $2 FROM deposits
           WHERE sub_address=$3`,
					vLog.TxHash.Hex(), amt.String(), sa.Hex(),
				)
				if err != nil {
					log.Println("写入 payment 失败:", err)
				} else {
					log.Printf("记录充值: sub_address=%s amount=%s tx=%s\n", sa.Hex(), amt.String(), vLog.TxHash.Hex())
				}
				break
			}
		}
	}
}
