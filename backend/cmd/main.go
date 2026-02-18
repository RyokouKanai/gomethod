package main

import (
	"log"
	"os"

	"github.com/RyokouKanai/gomethod/database"
	"github.com/RyokouKanai/gomethod/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// データベース接続
	database.Connect()

	// Gin ルーター設定
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// LINE Webhook
	r.POST("/callback", handler.WebhookHandler)

	// バッチ実行エンドポイント（Cloud Scheduler から OIDC 認証で呼び出し）
	batchGroup := r.Group("/batch")
	{
		batchGroup.POST("/:name", handler.BatchHandler)
	}

	// ポート設定（Cloud Run は PORT 環境変数を使用）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
