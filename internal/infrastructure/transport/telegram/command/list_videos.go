package command

import (
	"context"
	"fmt"
	"html"
	"strings"
)

func (r *Router) listVideos(ctx context.Context, chatID int64) {
	r.log.Debug("start listVideos")

	res, err := r.listVideosQuery.Run(ctx)
	if err != nil {
		r.send(chatID, fmt.Sprintf("Ошибка: %v", err))
		return
	}

	if len(res.Items) == 0 {
		r.send(chatID, "Список видео пуст")
		return
	}

	var sb strings.Builder
	sb.WriteString("<b>📋 Список видео</b>\n\n")

	for _, v := range res.Items {
		sb.WriteString(fmt.Sprintf(
			"📺 %s | 👁 %s | 👍 %s | 💬 %s\n🎬 <a href=\"%s\">%s</a>\n\n",
			platformShort(v.URL),
			humanize(v.Views),
			humanize(v.Likes),
			humanize(v.Comments),
			html.EscapeString(v.URL),
			html.EscapeString(trim(v.Title, 40)),
		))
	}

	r.send(chatID, sb.String())
}


func platformShort(url string) string {
	switch {
	case strings.Contains(url, "youtube"):
		return "y"
	case strings.Contains(url, "tiktok"):
		return "t"
	case strings.Contains(url, "instagram"):
		return "i"
	default:
		return "WEB"
	}
}

func humanize(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fm", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fk", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

func trim(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}