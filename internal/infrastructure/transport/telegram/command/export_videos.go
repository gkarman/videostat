package command

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/infrastructure/exporter/excel"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (r *Router) exportVideos(ctx context.Context, chatID int64) {
	r.log.Debug("start export videos")

	res, err := r.listVideosQuery.Run(ctx)
	if err != nil {
		r.send(chatID, fmt.Sprintf("Ошибка: %v", err))
		return
	}

	f, err := excel.BuildVideosSheet(res.Items)
	if err != nil {
		r.send(chatID, "Ошибка генерации Excel")
		return
	}

	defer func() {
		_ = f.Close()
	}()

	var buf *bytes.Buffer
	buf, err = f.WriteToBuffer()
	if err != nil {
		r.send(chatID, "Ошибка записи Excel файла")
		return
	}

	file := tgbotapi.FileBytes{
		Name:  "videos.xlsx",
		Bytes: buf.Bytes(),
	}

	msg := tgbotapi.NewDocument(chatID, file)
	msg.Caption = "📊 Экспорт видео"

	_ = r.sender.Send(msg)
}
