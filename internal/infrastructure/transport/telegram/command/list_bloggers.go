package command

import (
	"context"
	"fmt"
	"html"
	"strings"
)

func (r *Router) listBloggers(ctx context.Context, chatID int64) {
	r.log.Debug("start listBloggers")

	res, err := r.listBloggersQuery.Run(ctx)
	if err != nil {
		r.send(chatID, fmt.Sprintf("Ошибка: %v", err))
		return
	}

	if len(res.Items) == 0 {
		r.send(chatID, "Список блогеров пуст")
		return
	}

	var sb strings.Builder
	sb.WriteString("<b>📋 Список блогеров</b>\n\n")

	for i, bl := range res.Items {
		sb.WriteString(
			fmt.Sprintf(
				"%d. %s — <a href=\"%s\">\"%s\"</a>\n",
				i+1,
				platformIcon(bl.Platform),
				html.EscapeString(bl.URL),
				html.EscapeString(bl.URL),
			),
		)
	}

	r.send(chatID, sb.String())
}

func platformIcon(p string) string {
	switch strings.ToLower(p) {
	case "youtube":
		return "YouTube"
	case "tiktok":
		return "TikTok"
	case "instagram":
		return "Instagram"
	default:
		return strings.Title(p)
	}
}
