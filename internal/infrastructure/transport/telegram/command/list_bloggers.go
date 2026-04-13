package command

import (
	"context"
	"fmt"
	"strings"
)

func (r *Router) listBloggers(ctx context.Context, chatID int64) {
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
	sb.WriteString("<b>Список блогеров:</b>\n\n")

	for i, bl := range res.Items {
		sb.WriteString(
			fmt.Sprintf("%d. <a href='%s'>%s</a> (%s)\n",
				i+1,
				bl.URL,
				bl.URL,
				bl.Platform,
			),
		)
	}

	r.send(chatID, sb.String())
}
