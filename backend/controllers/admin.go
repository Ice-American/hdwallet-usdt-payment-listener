package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/Ice-American/hdwallet-usdt-payment-listener/backend/models"
	"github.com/Ice-American/hdwallet-usdt-payment-listener/backend/wallets"
	log "github.com/sirupsen/logrus"
)

// UserInfo 包含 userId 与其子存款次数
type UserInfo struct {
	UserID string `json:"user_id"`
	Count  int    `json:"count"`
}

// DepositInfo 包含单个子地址与索引
type DepositInfo struct {
	Index      int    `json:"index"`
	SubAddress string `json:"sub_address"`
}

// PrivateKeyResponse 返回单个子账户的私钥
type PrivateKeyResponse struct {
	PrivateKey string `json:"private_key"`
}

// HandleListUsers 列出所有有过充值请求的 userId
func HandleListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := models.DB.Query("SELECT user_id, COUNT(*) FROM deposits GROUP BY user_id")
	if err != nil {
		log.WithError(err).Error("查询用户列表失败")
		http.Error(w, "list users failed", 500)
		return
	}
	defer rows.Close()

	var list []UserInfo
	for rows.Next() {
		var u UserInfo
		rows.Scan(&u.UserID, &u.Count)
		list = append(list, u)
	}
	json.NewEncoder(w).Encode(list)
}

// HandleListDeposits 列出指定 userId 的所有子地址
func HandleListDeposits(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "missing userId", 400)
		return
	}
	rows, err := models.DB.Query(
		"SELECT derivation_index, sub_address FROM deposits WHERE user_id=$1 ORDER BY derivation_index",
		userId,
	)
	if err != nil {
		log.WithError(err).Error("查询 deposits 失败")
		http.Error(w, "list deposits failed", 500)
		return
	}
	defer rows.Close()

	var deps []DepositInfo
	for rows.Next() {
		var d DepositInfo
		rows.Scan(&d.Index, &d.SubAddress)
		deps = append(deps, d)
	}
	json.NewEncoder(w).Encode(deps)
}

// HandleGetPrivateKey 根据 userId 和 index 派生并返回私钥
func HandleGetPrivateKey(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	idxStr := r.URL.Query().Get("index")
	if userId == "" || idxStr == "" {
		http.Error(w, "missing parameters", 400)
		return
	}
	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		http.Error(w, "invalid index", 400)
		return
	}
	mnemonic := os.Getenv("MNEMONIC")
	priv, err := wallets.DerivePrivateKey(mnemonic, uint32(idx))
	if err != nil {
		log.WithError(err).Error("派生私钥失败")
		http.Error(w, "derive key failed", 500)
		return
	}
	json.NewEncoder(w).Encode(PrivateKeyResponse{PrivateKey: priv})
}
