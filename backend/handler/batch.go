package handler

import (
	"net/http"

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
