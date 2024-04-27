package main

import (
	"log"
	"time"

	"apu-ptt/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	// Routes
	r.POST("/api/import/dttot", handler.ImportTeroristHandler)
	r.POST("/api/import/ppspm", handler.ImportPpspmtHandler)
	r.POST("/api/import/transaksi", handler.ImportRegisteredMemberTransactionsHandler)
	r.GET("/api/unregistered_member_transactions", handler.ImportUnregisteredMemberTransactionsHandler)
	r.GET("/api/registered_transactions", handler.ImportRegisteredMemberTransactionsHandler)
	r.GET("/api/tkm1", handler.CustomersWithMultipleTransactionsHandler)
	r.GET("/api/alltransaction", handler.ListAllTransactions)
	// r.GET("/api/tkm", handler.TransactionsExceedingTwoTimesGroupedHandler)

	// Run server
	if err := r.Run(":9080"); err != nil {
		log.Fatal(err)
	}
}
