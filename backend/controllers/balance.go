// controllers/balance.go
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Ice-American/hdwallet-usdt-payment-listener/backend/models"
)

type BalanceResponse struct {
	TotalDeposited string `json:"total_deposited"`
	TotalConfirmed string `json:"total_confirmed"`
}

func HandleBalance(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	// 累计所有请求金额（可选：基于 payments 表）
	var sum string
	models.DB.QueryRow(
		`SELECT COALESCE(SUM(amount),0) FROM payments
     WHERE deposit_id IN (
       SELECT id FROM deposits WHERE user_id=$1
     )`, userId,
	).Scan(&sum)
	json.NewEncoder(w).Encode(BalanceResponse{TotalDeposited: sum})
}
