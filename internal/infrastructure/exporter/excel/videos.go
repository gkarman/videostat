package excel

import (
	"fmt"
	"strings"

	"github.com/gkarman/demo/internal/application/blogger/query/view"
	"github.com/xuri/excelize/v2"
)

func BuildVideosSheet(items []*view.Video) (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := "Videos"

	err := f.SetSheetName("Sheet1", sheet)
	if err != nil {
		return nil, fmt.Errorf("set sheet name: %w", err)
	}

	headers := []string{
		"Platform",
		"Title",
		"Views",
		"Likes",
		"Comments",
		"PublishedAt",
		"CreatedAt",
		"URL",
		"BloggerURL",
	}

	// header
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		err = f.SetCellValue(sheet, cell, h)
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", i, err)
		}
	}

	for i, v := range items {
		row := i + 2

		err = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), platformShort(v.URL))
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", row, err)
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Title)
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", row, err)
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Views)
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", row, err)
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.Likes)
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", row, err)
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("E%d", row), v.Comments)
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", row, err)
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("F%d", row), v.PublishedAt.Format("2006-01-02 15:04"))
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", row, err)
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("G%d", row), v.CreatedAt.Format("2006-01-02 15:04"))
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", row, err)
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("H%d", row), v.URL)
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", row, err)
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("I%d", row), v.BloggerURL)
		if err != nil {
			return nil, fmt.Errorf("set cell %d: %w", row, err)
		}
	}

	return f, nil
}

func platformShort(url string) string {
	switch {
	case strings.Contains(url, "youtube"):
		return "YouTube"
	case strings.Contains(url, "tiktok"):
		return "TikTok"
	case strings.Contains(url, "instagram"):
		return "Instagram"
	default:
		return "Web"
	}
}
