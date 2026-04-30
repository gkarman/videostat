package excel

import (
	"fmt"
	"strings"

	"github.com/gkarman/demo/internal/application/blogger/query/view"
	"github.com/xuri/excelize/v2"
)

func BuildVideosWorkbook(items []*view.Video) (*excelize.File, error) {
	f := excelize.NewFile()

	grouped := groupByPlatform(items)

	first := true
	for platform, videos := range grouped {
		sheet := strings.ToUpper(platform[:1]) + platform[1:]

		if first {
			if err := f.SetSheetName("Sheet1", sheet); err != nil {
				return nil, fmt.Errorf("set sheet name: %w", err)
			}
			first = false
		} else {
			if _, err := f.NewSheet(sheet); err != nil {
				return nil, fmt.Errorf("new sheet: %w", err)
			}
		}

		if err := writeVideosSheet(f, sheet, videos); err != nil {
			return nil, fmt.Errorf("write sheet %s: %w", sheet, err)
		}
	}

	return f, nil
}

func groupByPlatform(items []*view.Video) map[string][]*view.Video {
	res := make(map[string][]*view.Video)
	for _, v := range items {
		p := strings.ToLower(v.Platform)
		res[p] = append(res[p], v)
	}
	return res
}

var headers = []string{
	"BloggerURL",
	"Title",
	"URL",
	"PublishedAt",
	"CreatedAt",
	"Views",
	"Likes",
	"Comments",
	"Viral",
	"Relevant",
	"Status",
	"ErrorStage",
	"ErrorMessage",
	"AnalysisProvider",
	"AnalysisData",
	"LLMProvider",
	"LLMModel",
	"Prompt",
}

func writeVideosSheet(f *excelize.File, sheet string, videos []*view.Video) error {
	lastCol, _ := excelize.ColumnNumberToName(len(headers))

	if err := f.SetColWidth(sheet, "A", "A", 45); err != nil {
		return fmt.Errorf("set col width A: %w", err)
	}
	if err := f.SetColWidth(sheet, "B", "J", 20); err != nil {
		return fmt.Errorf("set col width B-J: %w", err)
	}
	// Status, ErrorStage, ErrorMessage
	if err := f.SetColWidth(sheet, "K", "M", 20); err != nil {
		return fmt.Errorf("set col width K-M: %w", err)
	}
	// AnalysisProvider
	if err := f.SetColWidth(sheet, "N", "N", 18); err != nil {
		return fmt.Errorf("set col width N: %w", err)
	}
	if err := f.SetColWidth(sheet, "O", "R", 20); err != nil {
		return fmt.Errorf("set col width O-R: %w", err)
	}

	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: false},
	})
	if err != nil {
		return fmt.Errorf("new style: %w", err)
	}
	if err := f.SetColStyle(sheet, "A:"+lastCol, style); err != nil {
		return fmt.Errorf("set col style: %w", err)
	}

	for i, h := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return fmt.Errorf("header coord: %w", err)
		}
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return fmt.Errorf("set header %s: %w", cell, err)
		}
	}

	set := func(row, col int, value any) error {
		cell, err := excelize.CoordinatesToCellName(col, row)
		if err != nil {
			return err
		}
		return f.SetCellValue(sheet, cell, value)
	}

	deref := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	for i, v := range videos {
		row := i + 2

		viral := ""
		if v.Viral {
			viral = "+"
		}
		relevant := "-"
		if v.IsRelevant {
			relevant = "+"
		}

		vals := []any{
			v.BloggerURL,
			v.Title,
			v.URL,
			v.PublishedAt.Format("2006-01-02"),
			v.CreatedAt.Format("2006-01-02"),
			v.Views,
			v.Likes,
			v.Comments,
			viral,
			relevant,
			v.Status,
			deref(v.ErrorStage),
			deref(v.ErrorMessage),
			deref(v.AnalysisProvider),
			deref(v.RawPayload),
			deref(v.LLMProvider),
			deref(v.LLMModel),
			deref(v.Prompt),
		}

		for col, val := range vals {
			if err := set(row, col+1, val); err != nil {
				return err
			}
		}
	}

	for i := 1; i <= len(videos)+1; i++ {
		if err := f.SetRowHeight(sheet, i, 18); err != nil {
			return err
		}
	}

	return nil
}
