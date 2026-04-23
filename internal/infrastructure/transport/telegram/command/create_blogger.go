package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (r *Router) askPlatform(chatID int64) {
	r.sendWithKeyboard(
		chatID,
		"Выбери платформу:",
		r.ui.PlatformKeyboard(),
	)
}

func (r *Router) handleFSM(ctx context.Context, msg *tgbotapi.Message) bool {
	st, ok := r.state.Get(msg.Chat.ID)
	if !ok {
		return false
	}

	if st.WaitingVideoURL {
		videoURL := msg.Text
		resp, err := r.startProcessVideo.Run(ctx, reqdto.StartProcessVideo{
			URL: videoURL,
		})
		if err != nil {
			r.log.Error("telegram request failed", "err", err, "url", videoURL)
			r.send(msg.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
			return true
		}

		r.state.Clear(msg.Chat.ID)
		r.send(msg.Chat.ID, resp.Message)

		return true
	}

	if st.WaitingURL {
		resp, err := r.createBlogger.Run(ctx, reqdto.CreateBlogger{
			URL:          msg.Text,
			PlatformName: st.PlatformName,
		})
		if err != nil {
			r.send(msg.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
			return true
		}

		r.state.Clear(msg.Chat.ID)
		r.send(msg.Chat.ID, fmt.Sprintf("Блогер создан ✅\nID: <code>%s</code>", resp.ID))
		return true
	}

	return false
}

func (r *Router) handleCallback(ctx context.Context, q *tgbotapi.CallbackQuery) {
	switch {
	case q.Data == "create_blogger":
		r.askPlatform(q.Message.Chat.ID)

	case strings.HasPrefix(q.Data, "platform_"):
		platform := strings.TrimPrefix(q.Data, "platform_")
		r.state.Set(q.Message.Chat.ID, &userState{
			PlatformName: platform,
			WaitingURL:   true,
		})
		r.send(q.Message.Chat.ID, "Пришли ссылку на блогера")

	case q.Data == "list_bloggers":
		r.listBloggers(ctx, q.Message.Chat.ID)

	case q.Data == "list_videos":
		r.listVideos(ctx, q.Message.Chat.ID)

	case q.Data == "export_videos":
		r.exportVideos(ctx, q.Message.Chat.ID)

	case q.Data == "start_process_video":
		r.state.Set(q.Message.Chat.ID, &userState{
			WaitingVideoURL: true,
		})
		r.send(q.Message.Chat.ID, "Пришли ссылку на видео")
	}
}
