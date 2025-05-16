package main

import (
	"net/http"
	"os"

	"github.com/Ice-American/hdwallet-usdt-payment-listener/backend/controllers"
	"github.com/Ice-American/hdwallet-usdt-payment-listener/backend/listener"
	"github.com/Ice-American/hdwallet-usdt-payment-listener/backend/models"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	// 输出到 stdout
	log.SetOutput(os.Stdout)
	// 日志等级
	log.SetLevel(log.InfoLevel)

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: no .env file found")
	}
	models.InitDB()
	go listener.ListenUSDT()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/admin/users", controllers.HandleListUsers)
	mux.HandleFunc("/api/admin/deposits", controllers.HandleListDeposits)
	mux.HandleFunc("/api/admin/private-key", controllers.HandleGetPrivateKey)

	// 默认允许所有
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // 允许所有域
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	handler := c.Handler(mux)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
