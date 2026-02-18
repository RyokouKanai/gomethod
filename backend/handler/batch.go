package handler

import (
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/RyokouKanai/gomethod/batch"
	"github.com/gin-gonic/gin"
)

// バッチ名とバッチ関数のマッピング
var batchRegistry = map[string]func(){
	"send_daily_g_message":       batch.SendDailyGMessage,
	"send_weekly_g_message":      batch.SendWeeklyGMessage,
	"send_weekly_blog_g_message": batch.SendWeeklyBlogGMessage,
	"send_experience_g_message":  batch.SendExperienceGMessage,
	"send_moon_message_today":    batch.SendMoonMessageToday,
	"send_moon_message_tomorrow": batch.SendMoonMessageTomorrow,
	"send_notice":                batch.SendNotice,
}

// BatchAuthMiddleware validates the batch request using a shared secret token.
// Cloud Scheduler sends this token in the Authorization header.
func BatchAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := os.Getenv("BATCH_AUTH_TOKEN")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "batch auth not configured"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		expected := "Bearer " + token

		if subtle.ConstantTimeCompare([]byte(authHeader), []byte(expected)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}

// BatchHandler executes a batch job by name.
// POST /batch/:name
func BatchHandler(c *gin.Context) {
	name := c.Param("name")

	fn, ok := batchRegistry[name]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown batch: " + name})
		return
	}

	// 非同期で実行（Cloud Schedulerのタイムアウトを避ける）
	go fn()

	c.JSON(http.StatusOK, gin.H{"status": "started", "batch": name})
}
