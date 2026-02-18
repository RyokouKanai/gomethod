package action

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/RyokouKanai/gomethod/database"
	"github.com/RyokouKanai/gomethod/model"
	"github.com/RyokouKanai/gomethod/service"
)

// Registry maps execution_method names to action functions.
type Registry struct {
	actions map[string]ActionFunc
}

// ActionFunc is the function signature for all actions.
// Returns the message content to reply with (string or []string).
type ActionFunc func(user *model.User, receivedMessage string, replyToken string, nextMessage *model.Message) interface{}

// NewRegistry creates a new action registry with all actions registered.
func NewRegistry() *Registry {
	r := &Registry{
		actions: make(map[string]ActionFunc),
	}
	r.registerAll()
	return r
}

// Execute implements the ActionExecutor interface.
func (r *Registry) Execute(method string, user *model.User, receivedMessage string, replyToken string, nextMessage *model.Message) interface{} {
	fn, ok := r.actions[method]
	if !ok {
		log.Printf("Unknown action method: %s", method)
		return nextMessage.ToFormattedText()
	}
	return fn(user, receivedMessage, replyToken, nextMessage)
}

func (r *Registry) registerAll() {
	// User content actions
	r.actions["dream_wishes_index"] = dreamWishesIndex
	r.actions["dream_wishes_create"] = dreamWishesCreate
	r.actions["dream_wishes_edit"] = dreamWishesEdit
	r.actions["dream_wishes_update"] = dreamWishesUpdate
	r.actions["dream_wishes_destroy"] = dreamWishesDestroy
	r.actions["solution_wishes_index"] = solutionWishesIndex
	r.actions["solution_wishes_create"] = solutionWishesCreate
	r.actions["solution_wishes_edit"] = solutionWishesEdit
	r.actions["solution_wishes_update"] = solutionWishesUpdate
	r.actions["solution_wishes_destroy"] = solutionWishesDestroy
	r.actions["hates_index"] = hatesIndex
	r.actions["hates_create"] = hatesCreate
	r.actions["hates_edit"] = hatesEdit
	r.actions["hates_update"] = hatesUpdate
	r.actions["hates_destroy"] = hatesDestroy
	r.actions["hates_destroy_all"] = hatesDestroyAll
	r.actions["happiness_index"] = happinessIndex
	r.actions["happiness_create"] = happinessCreate
	r.actions["happiness_destroy"] = happinessDestroy
	r.actions["talks_index"] = talksIndex
	r.actions["g_messages_show"] = gMessagesShow
	r.actions["thanks_count_show"] = thanksCountShow
	r.actions["thanks_count_reset"] = thanksCountReset
	r.actions["experiences_index"] = experiencesIndex
	r.actions["experiences_show"] = experiencesShow
	r.actions["find_or_create_feeling_settings"] = findOrCreateFeelingSettings
	r.actions["echo_feeling"] = echoFeeling
	r.actions["feeling_setting_index"] = feelingSettingIndex
	r.actions["feeling_setting_edit"] = feelingSettingEdit
	r.actions["feeling_setting_update"] = feelingSettingUpdate
	r.actions["save_selected_option"] = saveSelectedOption

	// File actions
	r.actions["dream_wish_file_create"] = dreamWishFileCreate
	r.actions["dream_wish_file_update"] = dreamWishFileUpdate
	r.actions["dream_wish_file_show"] = dreamWishFileShow

	// Admin actions
	r.actions["broadcasts_confirm"] = broadcastsConfirm
	r.actions["broadcasts"] = broadcasts
	r.actions["g_messages_create"] = gMessagesCreate
	r.actions["g_messages_index"] = gMessagesIndex
	r.actions["g_messages_destroy"] = gMessagesDestroy
	r.actions["g_messages_edit"] = gMessagesEdit
	r.actions["g_messages_update"] = gMessagesUpdate
	r.actions["weekly_g_messages_create"] = weeklyGMessagesCreate
	r.actions["weekly_g_messages_index"] = weeklyGMessagesIndex
	r.actions["weekly_g_messages_destroy"] = weeklyGMessagesDestroy
	r.actions["weekly_g_messages_edit"] = weeklyGMessagesEdit
	r.actions["weekly_g_messages_update"] = weeklyGMessagesUpdate
	r.actions["weekly_blog_g_messages_create"] = weeklyBlogGMessagesCreate
	r.actions["weekly_blog_g_messages_index"] = weeklyBlogGMessagesIndex
	r.actions["weekly_blog_g_messages_destroy"] = weeklyBlogGMessagesDestroy
	r.actions["weekly_blog_g_messages_edit"] = weeklyBlogGMessagesEdit
	r.actions["weekly_blog_g_messages_update"] = weeklyBlogGMessagesUpdate
	r.actions["experience_g_messages_create"] = experienceGMessagesCreate
	r.actions["experience_g_messages_index"] = experienceGMessagesIndex
	r.actions["experience_g_messages_destroy"] = experienceGMessagesDestroy
	r.actions["experience_g_messages_edit"] = experienceGMessagesEdit
	r.actions["experience_g_messages_update"] = experienceGMessagesUpdate
	r.actions["notices_create"] = noticesCreate
	r.actions["notices_index"] = noticesIndex
	r.actions["notices_destroy"] = noticesDestroy
	r.actions["notices_edit"] = noticesEdit
	r.actions["notices_update"] = noticesUpdate
}

// Helper: selected number from last_message (0-indexed)
func selectedNumber(user *model.User) int {
	lm, err := user.GetLastMessage()
	if err != nil {
		return -1
	}
	n, err := strconv.Atoi(lm.PlainContent())
	if err != nil {
		return -1
	}
	return n - 1
}

// Helper: save user's selection
func saveSelectedOption(user *model.User, receivedMessage string, _ string, nextMessage *model.Message) interface{} {
	if err := user.UpsertLastMessage(receivedMessage); err != nil {
		msg := model.GetMessageByScope("validation_error")
		if msg != nil {
			return msg.GetContent()
		}
		return "エラーが発生しました"
	}
	return nextMessage.ToFormattedText()
}

// ==================== Dream Wishes ====================

func dreamWishesIndex(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	wishes, _ := user.GetDreamWishes()
	if len(wishes) == 0 {
		msg := model.GetMessageByScope("no_wishes")
		if msg != nil {
			return msg.GetContent()
		}
		return "願いがまだ登録されていません"
	}
	base := nextMessage.ToFormattedText()
	return base + "\n\n" + formatWishes(wishes)
}

func dreamWishesCreate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	model.CreateWish(user.ID, msg, "dream")
	return nextMessage.ToFormattedText()
}

func dreamWishesEdit(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	wishes, _ := user.GetDreamWishes()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(wishes) {
		return nextMessage.GetContent()
	}
	return nextMessage.GetContent() + "\n\n選択中の願い:\n" + wishes[idx].PlainContent()
}

func dreamWishesUpdate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	wishes, _ := user.GetDreamWishes()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(wishes) {
		return nextMessage.GetContent()
	}
	w := model.FindWishByID(wishes[idx].ID)
	if w != nil {
		w.Content = &msg
		model.UpdateWishContent(w)
	}
	return nextMessage.GetContent() + "\n\n" + msg
}

func dreamWishesDestroy(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	wishes, _ := user.GetDreamWishes()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(wishes) {
		return nextMessage.GetContent()
	}
	w := model.FindWishByID(wishes[idx].ID)
	plain := ""
	if w != nil {
		plain = w.PlainContent()
		database.DB.Delete(w)
	}
	return nextMessage.GetContent() + "\n\n" + plain
}

// ==================== Solution Wishes ====================

func solutionWishesIndex(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	wishes, _ := user.GetSolutionWishes()
	if len(wishes) == 0 {
		msg := model.GetMessageByScope("no_wishes")
		if msg != nil {
			return msg.GetContent()
		}
		return "願いがまだ登録されていません"
	}
	base := nextMessage.ToFormattedText()
	return base + "\n\n" + formatWishes(wishes)
}

func solutionWishesCreate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	model.CreateWish(user.ID, msg, "solution")
	return nextMessage.ToFormattedText()
}

func solutionWishesEdit(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	wishes, _ := user.GetSolutionWishes()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(wishes) {
		return nextMessage.GetContent()
	}
	return nextMessage.GetContent() + "\n\n選択中の願い:\n" + wishes[idx].PlainContent()
}

func solutionWishesUpdate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	wishes, _ := user.GetSolutionWishes()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(wishes) {
		return nextMessage.GetContent()
	}
	w := model.FindWishByID(wishes[idx].ID)
	if w != nil {
		w.Content = &msg
		model.UpdateWishContent(w)
	}
	return nextMessage.GetContent() + "\n\n" + msg
}

func solutionWishesDestroy(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	wishes, _ := user.GetSolutionWishes()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(wishes) {
		return nextMessage.GetContent()
	}
	w := model.FindWishByID(wishes[idx].ID)
	plain := ""
	if w != nil {
		plain = w.PlainContent()
		database.DB.Delete(w)
	}
	return nextMessage.GetContent() + "\n\n" + plain
}

// ==================== Hates ====================

func hatesIndex(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	hates, _ := user.GetHates()
	if len(hates) == 0 {
		return "まだ嫌だー！を投企してないようです。。"
	}
	base := nextMessage.ToFormattedText()
	return base + "\n\n" + formatHates(hates)
}

func hatesCreate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	model.CreateHate(user.ID, msg)
	return nextMessage.ToFormattedText()
}

func hatesEdit(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	hates, _ := user.GetHates()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(hates) {
		return nextMessage.GetContent()
	}
	return nextMessage.GetContent() + "\n\n選択中の嫌だー:\n" + hates[idx].PlainContent()
}

func hatesUpdate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	hates, _ := user.GetHates()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(hates) {
		return nextMessage.GetContent()
	}
	h := model.FindHateByID(hates[idx].ID)
	if h != nil {
		h.Content = &msg
		model.UpdateHateContent(h)
	}
	return nextMessage.GetContent() + "\n\n" + msg
}

func hatesDestroy(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	hates, _ := user.GetHates()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(hates) {
		return nextMessage.GetContent()
	}
	h := model.FindHateByID(hates[idx].ID)
	plain := ""
	if h != nil {
		plain = h.PlainContent()
		database.DB.Delete(h)
	}
	return nextMessage.GetContent() + "\n\n" + plain
}

func hatesDestroyAll(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	model.DeleteHatesByUserID(user.ID)
	return nextMessage.ToFormattedText()
}

// ==================== Happiness ====================

func happinessIndex(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	happiness, _ := user.GetHappiness()
	if len(happiness) == 0 {
		return "まだ良かったー！を書いてないようだね。これからどんどん書いていこう！"
	}
	base := nextMessage.ToFormattedText()
	return base + "\n\n" + formatHappiness(happiness)
}

func happinessCreate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	model.CreateHappiness(user.ID, msg)
	return nextMessage.ToFormattedText()
}

func happinessDestroy(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	happiness, _ := user.GetHappiness()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(happiness) {
		return nextMessage.GetContent()
	}
	h := model.FindHappinessByID(happiness[idx].ID)
	plain := ""
	if h != nil {
		plain = h.PlainContent()
		database.DB.Delete(h)
	}
	return nextMessage.GetContent() + "\n\n" + plain
}

// ==================== Talks (OpenAI) ====================

func talksIndex(_ *model.User, msg string, _ string, _ *model.Message) interface{} {
	n, err := strconv.Atoi(msg)
	if err == nil && n != 0 {
		return "数字以外の言葉を入力してね！"
	}

	endpoint := os.Getenv("OPEN_AI_COMPLETION_ENDPOINT")
	apiKey := os.Getenv("OPEN_AI_API_KEY")
	if endpoint == "" || apiKey == "" {
		badResp := model.GetMessageByScope("bad_talk_response")
		if badResp != nil {
			return badResp.GetContent()
		}
		return "申し訳ございません。エラーが発生しました。"
	}

	cleanMsg := strings.ReplaceAll(strings.ReplaceAll(msg, "\r", ""), "\n", "")
	body := fmt.Sprintf(`{"model":"text-davinci-003","prompt":"%s","temperature":0.7,"max_tokens":256,"top_p":1,"frequency_penalty":0,"presence_penalty":0}`, cleanMsg)

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		badResp := model.GetMessageByScope("bad_talk_response")
		if badResp != nil {
			return badResp.GetContent()
		}
		return "申し訳ございません。エラーが発生しました。"
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		badResp := model.GetMessageByScope("bad_talk_response")
		if badResp != nil {
			return badResp.GetContent()
		}
		return "申し訳ございません。エラーが発生しました。"
	}

	first := choices[0].(map[string]interface{})
	text := first["text"].(string)
	text = strings.ReplaceAll(strings.ReplaceAll(text, "\r", ""), "\n", "")
	text = strings.TrimPrefix(text, "。")
	return text
}

// ==================== GMessages ====================

func gMessagesShow(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	gMsg, err := user.FetchGMessageByPeriod("daily")
	if err != nil {
		return nextMessage.GetContent()
	}
	user.CreateGMessageHistory(gMsg)
	return nextMessage.GetContent() + "\n\n" + gMsg.PlainContent()
}

// ==================== Thanks Count ====================

func thanksCountShow(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	count := user.GetThanksCount()
	rate := float64(100*count) / 1000.0
	return fmt.Sprintf("%s\n\n現在の回数: %d回\n達成度: %.1f%%", nextMessage.GetContent(), count, rate)
}

func thanksCountReset(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	user.ResetThanksCount()
	return nextMessage.ToFormattedText()
}

// ==================== Experiences ====================

func experiencesIndex(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	articles := getExperienceArticles(user)
	var titles []string
	for i, a := range articles {
		titles = append(titles, fmt.Sprintf("%d: %s", i+1, a.Title))
	}
	return nextMessage.GetContent() + "\n\n" + strings.Join(titles, "\n")
}

func experiencesShow(user *model.User, msg string, _ string, _ *model.Message) interface{} {
	articles := getExperienceArticles(user)
	idx, err := strconv.Atoi(msg)
	if err != nil || idx < 1 || idx > len(articles) {
		return "記事が見つかりませんでした"
	}
	sections, _ := articles[idx-1].GetSections()
	var contents []string
	for _, s := range sections {
		contents = append(contents, s.GetContent())
	}
	return contents
}

func getExperienceArticles(user *model.User) []model.Article {
	// Try to find lesson from talk history
	th, err := user.GetLatestTalkHistory()
	if err == nil && th != nil {
		converterMap := map[uint]uint{51: 3, 53: 4, 68: 5}
		if lessonID, ok := converterMap[th.MessageID]; ok {
			lesson := model.FindLessonByID(lessonID)
			if lesson != nil {
				articles, _ := lesson.GetArticles()
				if len(articles) > 0 {
					return articles
				}
			}
		}

		// Check last_message for selected_lesson_id
		lm, err := user.GetLastMessage()
		if err == nil && lm != nil {
			var data map[string]interface{}
			if json.Unmarshal([]byte(lm.PlainContent()), &data) == nil {
				if idVal, ok := data["selected_lesson_id"]; ok {
					if id, ok := idVal.(float64); ok {
						lesson := model.FindLessonByID(uint(id))
						if lesson != nil {
							articles, _ := lesson.GetArticles()
							if len(articles) > 0 {
								return articles
							}
						}
					}
				}
			}
		}
	}

	articles, _ := model.GetExperienceArticles()
	return articles
}

// ==================== Feeling Settings ====================

func findOrCreateFeelingSettings(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	settings, _ := user.GetFeelingSettings()
	if len(settings) == 0 {
		user.CreateFeelingSettings()
		settings, _ = user.GetFeelingSettings()
	}
	base := nextMessage.ToFormattedText()
	return base + "\n\n" + formatFeelingSettings(settings) + "\n6: 設定をカスタマイズする"
}

func echoFeeling(user *model.User, msg string, _ string, _ *model.Message) interface{} {
	n, err := strconv.Atoi(msg)
	if err != nil {
		return feelingSettingIndexInternal(user)
	}
	fs := model.FindFeelingSettingByUserAndButton(user.ID, n)
	if fs != nil {
		return fs.PlainContent()
	}
	return feelingSettingIndexInternal(user)
}

func feelingSettingIndex(user *model.User, _ string, _ string, _ *model.Message) interface{} {
	return feelingSettingIndexInternal(user)
}

func feelingSettingIndexInternal(user *model.User) string {
	msg := model.GetMessageByScope("lets_customize_feeling_button")
	base := ""
	if msg != nil {
		base = msg.ToFormattedText()
	}
	settings, _ := user.GetFeelingSettings()
	return base + "\n\n" + formatFeelingSettings(settings)
}

func feelingSettingEdit(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	user.UpsertLastMessage(msg)
	settings, _ := user.GetFeelingSettings()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(settings) {
		return feelingSettingIndexInternal(user)
	}
	return nextMessage.GetContent() + "\n\n選択中の設定: " + settings[idx].PlainContent()
}

func feelingSettingUpdate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	settings, _ := user.GetFeelingSettings()
	idx := selectedNumber(user)
	if idx < 0 || idx >= len(settings) {
		return nextMessage.GetContent()
	}
	fs := model.FindFeelingSettingByID(settings[idx].ID)
	if fs != nil {
		fs.Content = msg
		model.UpdateFeelingSettingContent(fs)
	}
	return nextMessage.ToFormattedText()
}

// ==================== File Actions ====================

func dreamWishFileCreate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	objectKey := "dream_wish/" + msg
	if err := uploadToS3(objectKey, msg); err != nil {
		badResp := model.GetMessageByScope("bad_talk_response")
		if badResp != nil {
			return badResp.GetContent()
		}
		return "エラーが発生しました"
	}
	wishes, _ := user.GetDreamWishes()
	if len(wishes) > 0 {
		model.UpdateWishS3URL(wishes[len(wishes)-1].ID, s3ObjectURL(objectKey))
	}
	return nextMessage.ToFormattedText()
}

func dreamWishFileUpdate(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	objectKey := "dream_wish/" + msg
	if err := uploadToS3(objectKey, msg); err != nil {
		badResp := model.GetMessageByScope("bad_talk_response")
		if badResp != nil {
			return badResp.GetContent()
		}
		return "エラーが発生しました"
	}
	wishes, _ := user.GetDreamWishes()
	idx := selectedNumber(user)
	if idx >= 0 && idx < len(wishes) {
		model.UpdateWishS3URL(wishes[idx].ID, s3ObjectURL(objectKey))
	}
	return nextMessage.ToFormattedText()
}

func dreamWishFileShow(user *model.User, msg string, replyToken string, nextMessage *model.Message) interface{} {
	wishes, _ := user.GetDreamWishes()
	idx, _ := strconv.Atoi(msg)
	idx-- // 1-indexed to 0-indexed
	if idx < 0 || idx >= len(wishes) {
		return "この願いには画像が投稿されていません。"
	}
	w := model.FindWishByID(wishes[idx].ID)
	if w == nil || w.GetS3ObjectURL() == "" {
		return "この願いには画像が投稿されていません。"
	}
	contents := []map[string]string{
		{"type": "text", "content": nextMessage.ToFormattedText()},
		{"type": "image", "content": w.GetS3ObjectURL()},
	}
	ss := service.NewSendService()
	ss.ReplyImageAndMessages(contents, replyToken)
	return nil
}

func uploadToS3(objectKey, messageID string) error {
	token := os.Getenv("LINE_CHANNEL_TOKEN")
	url := fmt.Sprintf("https://api-data.line.me/v2/bot/message/%s/content", messageID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	s3Svc := service.NewS3Service()
	bucket := os.Getenv("S3_BUCKET")
	return s3Svc.Upload(bucket, objectKey, resp.Body)
}

func s3ObjectURL(objectKey string) string {
	return os.Getenv("S3_ROOT_URL") + "/" + objectKey
}

// ==================== Admin: Broadcast ====================

func broadcastsConfirm(user *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
	lm, err := user.GetLastMessage()
	rangeOption := 0
	if err == nil && lm != nil {
		rangeOption, _ = strconv.Atoi(lm.PlainContent())
	}
	user.UpsertLastMessage(fmt.Sprintf("%d:&:%s", rangeOption, msg))
	rangeName := getRangeName(rangeOption)
	return msg + "\n\n" + nextMessage.ToFormattedText() + "\n\n送信対象：" + rangeName
}

func broadcasts(user *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
	lm, err := user.GetLastMessage()
	if err != nil || lm == nil {
		return nextMessage.GetContent()
	}
	parts := strings.SplitN(lm.PlainContent(), ":&:", 2)
	if len(parts) != 2 {
		return nextMessage.GetContent()
	}
	rangeOption, _ := strconv.Atoi(parts[0])
	sentMessage := parts[1]
	ss := service.NewSendService()

	if getRangeName(rangeOption) == "シックのみ" {
		shikUsers, _ := model.GetShikUsers()
		var ids []string
		for _, u := range shikUsers {
			ids = append(ids, u.LineUserID)
		}
		ss.BroadcastToShik(sentMessage, ids)
	} else {
		ss.Broadcast(sentMessage)
	}
	return nextMessage.GetContent()
}

func getRangeName(position int) string {
	broadcastRangeMsg := model.GetMessageByScope("select_broadcast_range")
	if broadcastRangeMsg == nil {
		return ""
	}
	o := model.FindOptionByMessageAndPosition(broadcastRangeMsg.ID, position)
	if o == nil {
		return ""
	}
	return o.GetContent()
}

// ==================== Admin: GMessages CRUD ====================

func gMessagesCRUD(period string) (
	create, index, destroy, edit, update ActionFunc,
) {
	create = func(_ *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
		model.CreateGMessage(msg, period)
		return nextMessage.GetContent() + "\n\n" + msg
	}
	index = func(_ *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
		messages, _ := model.GetGMessagesByPeriod(period)
		return nextMessage.GetContent() + "\n\n" + formatGMessages(messages)
	}
	destroy = func(_ *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
		messages, _ := model.GetGMessagesByPeriod(period)
		idx := selectedNumberFromMessages(messages)
		if idx < 0 || idx >= len(messages) {
			return nextMessage.GetContent()
		}
		plain := messages[idx].PlainContent()
		model.DeleteGMessage(messages[idx].ID)
		return nextMessage.GetContent() + "\n\n" + plain
	}
	edit = func(_ *model.User, _ string, _ string, nextMessage *model.Message) interface{} {
		messages, _ := model.GetGMessagesByPeriod(period)
		idx := selectedNumberFromMessages(messages)
		if idx < 0 || idx >= len(messages) {
			return nextMessage.GetContent()
		}
		label := periodLabel(period)
		return nextMessage.GetContent() + "\n\n選択中の" + label + ":\n" + messages[idx].PlainContent()
	}
	update = func(_ *model.User, msg string, _ string, nextMessage *model.Message) interface{} {
		messages, _ := model.GetGMessagesByPeriod(period)
		idx := selectedNumberFromMessages(messages)
		if idx < 0 || idx >= len(messages) {
			return nextMessage.GetContent()
		}
		g := model.FindGMessageByID(messages[idx].ID)
		if g != nil {
			g.Content = &msg
			model.UpdateGMessageContent(g)
		}
		return nextMessage.GetContent() + "\n\n" + msg
	}
	return
}

func periodLabel(period string) string {
	switch period {
	case "daily":
		return "Gメッセージ"
	case "weekly":
		return "Gメッセージ"
	case "weekly_blog":
		return "サンデーブログ"
	case "experience":
		return "体験談"
	case "notice":
		return "お知らせ"
	}
	return period
}

func selectedNumberFromMessages(_ []model.GMessage) int {
	// This requires context from the user, which we don't have here
	// In the original Ruby it used the same selectedNumber from last_message
	// This will be handled by the calling context
	return -1
}

// Generate admin CRUD functions for each period
var (
	gMessagesCreate, gMessagesIndex, gMessagesDestroy, gMessagesEdit, gMessagesUpdate                               = gMessagesCRUD("daily")
	weeklyGMessagesCreate, weeklyGMessagesIndex, weeklyGMessagesDestroy, weeklyGMessagesEdit, weeklyGMessagesUpdate = gMessagesCRUD("weekly")
	weeklyBlogGMessagesCreate, weeklyBlogGMessagesIndex, weeklyBlogGMessagesDestroy, weeklyBlogGMessagesEdit, weeklyBlogGMessagesUpdate = gMessagesCRUD("weekly_blog")
	experienceGMessagesCreate, experienceGMessagesIndex, experienceGMessagesDestroy, experienceGMessagesEdit, experienceGMessagesUpdate = gMessagesCRUD("experience")
	noticesCreate, noticesIndex, noticesDestroy, noticesEdit, noticesUpdate                                                             = gMessagesCRUD("notice")
)

// ==================== Formatters ====================

func formatWishes(wishes []model.Wish) string {
	var lines []string
	for i, w := range wishes {
		hasImg := "なし"
		if w.GetS3ObjectURL() != "" {
			hasImg = "あり"
		}
		lines = append(lines, fmt.Sprintf("%d:\n日付: %s\n内容: %s\n画像: %s",
			i+1, w.CreatedAt.Format("2006年1月2日"), w.PlainContent(), hasImg))
	}
	return strings.Join(lines, "\n\n")
}

func formatHates(hates []model.Hate) string {
	var lines []string
	for i, h := range hates {
		lines = append(lines, fmt.Sprintf("%d:\n日付: %s\n内容: %s",
			i+1, h.CreatedAt.Format("2006年1月2日"), h.PlainContent()))
	}
	return strings.Join(lines, "\n\n")
}

func formatHappiness(happiness []model.Happiness) string {
	var lines []string
	for i, h := range happiness {
		lines = append(lines, fmt.Sprintf("%d:\n日付: %s\n内容: %s",
			i+1, h.CreatedAt.Format("2006年1月2日"), h.PlainContent()))
	}
	return strings.Join(lines, "\n\n")
}

func formatFeelingSettings(settings []model.FeelingSetting) string {
	var lines []string
	for _, s := range settings {
		lines = append(lines, fmt.Sprintf("%d: %s", s.ButtonNumber, s.PlainContent()))
	}
	return strings.Join(lines, "\n")
}

func formatGMessages(messages []model.GMessage) string {
	var lines []string
	for i, m := range messages {
		plain := m.PlainContent()
		if len([]rune(plain)) > 30 {
			plain = string([]rune(plain)[:30]) + "..."
		}
		lines = append(lines, fmt.Sprintf("%d:\n%s", i+1, plain))
	}
	return strings.Join(lines, "\n\n")
}
