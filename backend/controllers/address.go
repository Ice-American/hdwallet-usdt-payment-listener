package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/Ice-American/hdwallet-usdt-payment-listener/backend/models"
	"github.com/Ice-American/hdwallet-usdt-payment-listener/backend/wallets"
	log "github.com/sirupsen/logrus"
)

type AddressResponse struct {
	SubAddress string `json:"sub_address"`
	Index      int    `json:"index"`
	Error      string `json:"error,omitempty"`
}

// HandleAddress 每次请求都生成一个新的子地址，并写入 deposits
func HandleAddress(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "missing userId", http.StatusBadRequest)
		return
	}

	mnemonic := os.Getenv("MNEMONIC")
	if mnemonic == "" {
		log.Error("MNEMONIC 未配置")
		http.Error(w, "mnemonic not configured", http.StatusInternalServerError)
		return
	}

	// 计算本次派生索引：从 deposits 表取该用户已请求次数
	var idx int
	err := models.DB.QueryRow(
		"SELECT COUNT(*) FROM deposits WHERE user_id=$1",
		userId,
	).Scan(&idx)
	if err != nil {
		log.WithError(err).Error("COUNT deposits 失败")
		http.Error(w, "count deposits failed", http.StatusInternalServerError)
		return
	}

	log.Infof("为 userId=%s 第 %d 次派生子地址", userId, idx+1)
	subAddr, err := wallets.DeriveAddress(mnemonic, uint32(idx))
	if err != nil {
		log.WithError(err).Error("派生子地址失败")
		http.Error(w, "derive address failed", http.StatusInternalServerError)
		return
	}

	// 写入 deposits
	_, err = models.DB.Exec(
		`INSERT INTO deposits(user_id, derivation_index, sub_address)
     VALUES($1, $2, $3)`,
		userId, idx, subAddr,
	)
	if err != nil {
		log.WithError(err).Error("写入 deposits 失败")
		http.Error(w, "insert deposit failed", http.StatusInternalServerError)
		return
	}

	log.Infof("生成子地址 %s for userId=%s", subAddr, userId)
	json.NewEncoder(w).Encode(AddressResponse{
		SubAddress: subAddr,
		Index:      idx,
	})
}
