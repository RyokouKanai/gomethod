package batch

import (
	"log"
	"time"

	"github.com/RyokouKanai/gomethod/model"
	"github.com/RyokouKanai/gomethod/service"
)

// Base provides common batch functionality.
type Base struct {
	Name          string
	ExecutionTime float64
}

// IsDuplicate checks if this batch was already executed today.
func (b *Base) IsDuplicate() bool {
	return model.CheckBatchDuplicateExecution(b.Name)
}

// PrintResult logs the execution result.
func (b *Base) PrintResult() {
	log.Printf(`
--------------------------------------------------
DONE: %s
実行時間: %.2f秒
--------------------------------------------------
`, b.Name, b.ExecutionTime)
}

// Broadcast sends a message to all users.
func (b *Base) Broadcast(message string) {
	service.NewSendService().Broadcast(message)
}

// Unicast sends a message to a specific user.
func (b *Base) Unicast(lineUserID, message string) {
	service.NewSendService().Unicast(lineUserID, message)
}

// RunBatch executes a batch with timing and dedup checks.
func RunBatch(name string, fn func()) {
	base := &Base{Name: name}
	if base.IsDuplicate() {
		log.Printf("Batch %s already executed today, skipping", name)
		return
	}

	start := time.Now()
	fn()
	base.ExecutionTime = time.Since(start).Seconds()
	base.PrintResult()
}

// SendDailyGMessage sends the daily G message to all users.
func SendDailyGMessage() {
	RunBatch("SendDailyGMessage", func() {
		masterUser, err := model.GetMasterUser()
		if err != nil {
			log.Printf("Error getting master user: %v", err)
			return
		}
		gMsg, err := masterUser.FetchGMessageByPeriod("daily")
		if err != nil {
			log.Printf("Error fetching daily g_message: %v", err)
			return
		}
		masterUser.CreateGMessageHistory(gMsg)

		todaysMsg := model.GetMessageByScope("todays_g_message")
		content := ""
		if todaysMsg != nil {
			content = todaysMsg.GetContent()
		}
		message := content + "\n\n" + gMsg.PlainContent()

		base := &Base{}
		base.Broadcast(message)
	})
}

// SendWeeklyGMessage sends the weekly G message (Saturday video).
func SendWeeklyGMessage() {
	RunBatch("SendWeeklyGMessage", func() {
		masterUser, err := model.GetMasterUser()
		if err != nil {
			log.Printf("Error getting master user: %v", err)
			return
		}
		gMsg, err := masterUser.FetchGMessageByPeriod("weekly")
		if err != nil {
			log.Printf("Error fetching weekly g_message: %v", err)
			return
		}
		masterUser.CreateGMessageHistory(gMsg)

		todaysMsg := model.GetMessageByScope("todays_weekly_g_message")
		content := ""
		if todaysMsg != nil {
			content = todaysMsg.GetContent()
		}
		message := content + "\n\n" + gMsg.PlainContent()

		base := &Base{}
		base.Broadcast(message)
	})
}

// SendWeeklyBlogGMessage sends the weekly blog message (Sunday).
func SendWeeklyBlogGMessage() {
	RunBatch("SendWeeklyBlogGMessage", func() {
		masterUser, err := model.GetMasterUser()
		if err != nil {
			log.Printf("Error getting master user: %v", err)
			return
		}
		gMsg, err := masterUser.FetchGMessageByPeriod("weekly_blog")
		if err != nil {
			log.Printf("Error fetching weekly_blog g_message: %v", err)
			return
		}
		masterUser.CreateGMessageHistory(gMsg)

		todaysMsg := model.GetMessageByScope("todays_weekly_blog_g_message")
		content := ""
		if todaysMsg != nil {
			content = todaysMsg.GetContent()
		}
		message := content + "\n\n" + gMsg.PlainContent()

		base := &Base{}
		base.Broadcast(message)
	})
}

// SendExperienceGMessage sends experience messages (Tue/Thu).
func SendExperienceGMessage() {
	RunBatch("SendExperienceGMessage", func() {
		masterUser, err := model.GetMasterUser()
		if err != nil {
			log.Printf("Error getting master user: %v", err)
			return
		}
		gMsg, err := masterUser.FetchGMessageByPeriod("experience")
		if err != nil {
			log.Printf("Error fetching experience g_message: %v", err)
			return
		}
		masterUser.CreateGMessageHistory(gMsg)

		todaysMsg := model.GetMessageByScope("todays_experience_g_message")
		content := ""
		if todaysMsg != nil {
			content = todaysMsg.GetContent()
		}
		message := content + "\n\n" + gMsg.PlainContent()

		base := &Base{}
		base.Broadcast(message)
	})
}

// SendMoonMessageToday sends moon phase messages on the day of new/full moon.
func SendMoonMessageToday() {
	RunBatch("SendMoonMessageToday", func() {
		mp := model.GetMoonPhaseToday()
		if mp == nil {
			return
		}

		var scope string
		switch mp.Phase {
		case "new":
			scope = "new_moon_today"
		case "full":
			scope = "full_moon_today"
		default:
			return
		}

		msg := model.GetMessageByScope(scope)
		if msg == nil {
			return
		}

		base := &Base{}
		base.Broadcast(msg.GetContent())
	})
}

// SendMoonMessageTomorrow sends moon phase messages the day before new/full moon.
func SendMoonMessageTomorrow() {
	RunBatch("SendMoonMessageTomorrow", func() {
		mp := model.GetMoonPhaseTomorrow()
		if mp == nil {
			return
		}

		var scope string
		switch mp.Phase {
		case "new":
			scope = "new_moon_tomorrow"
		case "full":
			scope = "full_moon_tomorrow"
		default:
			return
		}

		msg := model.GetMessageByScope(scope)
		if msg == nil {
			return
		}

		base := &Base{}
		base.Broadcast(msg.GetContent())
	})
}

// SendNotice sends periodic notices (1st and 15th of month).
func SendNotice() {
	RunBatch("SendNotice", func() {
		masterUser, err := model.GetMasterUser()
		if err != nil {
			log.Printf("Error getting master user: %v", err)
			return
		}
		notice, err := masterUser.FetchGMessageByPeriod("notice")
		if err != nil {
			log.Printf("Error fetching notice: %v", err)
			return
		}
		masterUser.CreateGMessageHistory(notice)

		base := &Base{}
		base.Broadcast(notice.PlainContent())
	})
}
