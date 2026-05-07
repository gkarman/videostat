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
	"Ссылка",
	"Заголовок",
	"Ссылка",
	"Выложено",
	"Добавлено нам в БД",
	"Просмотры",
	"Лайки",
	"Комментарии",
	"Признак вирусности",
	"Признак релевантности заголовка",
	"Статус обработки",
	"Стадия на которой возника ошибка",
	"Ошибка",
	"Ресурс аналица видео",
	"Данные от ресурса анализа видео",
	"ИИ",
	"ИИ модель",
	"Бриф",
	"Сценарий",
	"Платформа генерации",
	"ID генерации",
	"Статус генерации",
	"Ошибка генерации",
	"Ссылка на S3",
}

type columnGroup struct {
	title    string
	fromCol  int
	toCol    int
}

var columnGroups = []columnGroup{
	{"Блогер", 1, 1},
	{"Данные о видео", 2, 8},
	{"Показатели качества контента", 9, 10},
	{"Анализ видео донора", 11, 15},
	{"Данные для генерации нового видео", 16, 19},
	{"Данные о новом видео", 20, 24},
}

func writeVideosSheet(f *excelize.File, sheet string, videos []*view.Video) error {
	lastCol, _ := excelize.ColumnNumberToName(len(headers))

	if err := f.SetColWidth(sheet, "A", "A", 60); err != nil {
		return fmt.Errorf("set col width A: %w", err)
	}
	if err := f.SetColWidth(sheet, "B", lastCol, 20); err != nil {
		return fmt.Errorf("set col width: %w", err)
	}

	dataStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: false},
	})
	if err != nil {
		return fmt.Errorf("new data style: %w", err)
	}
	if err := f.SetColStyle(sheet, "A:"+lastCol, dataStyle); err != nil {
		return fmt.Errorf("set col style: %w", err)
	}

	groupStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: false},
	})
	if err != nil {
		return fmt.Errorf("new group style: %w", err)
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"},
	})
	if err != nil {
		return fmt.Errorf("new header style: %w", err)
	}

	// row 1: group labels
	for _, g := range columnGroups {
		fromCell, _ := excelize.CoordinatesToCellName(g.fromCol, 1)
		toCell, _ := excelize.CoordinatesToCellName(g.toCol, 1)
		if g.fromCol != g.toCol {
			if err := f.MergeCell(sheet, fromCell, toCell); err != nil {
				return fmt.Errorf("merge group cells %s:%s: %w", fromCell, toCell, err)
			}
		}
		if err := f.SetCellValue(sheet, fromCell, g.title); err != nil {
			return fmt.Errorf("set group label %s: %w", fromCell, err)
		}
		if err := f.SetCellStyle(sheet, fromCell, toCell, groupStyle); err != nil {
			return fmt.Errorf("set group style %s: %w", fromCell, err)
		}
	}
	if err := f.SetRowHeight(sheet, 1, 30); err != nil {
		return fmt.Errorf("set group row height: %w", err)
	}

	// row 2: column headers
	for i, h := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 2)
		if err != nil {
			return fmt.Errorf("header coord: %w", err)
		}
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return fmt.Errorf("set header %s: %w", cell, err)
		}
		if err := f.SetCellStyle(sheet, cell, cell, headerStyle); err != nil {
			return fmt.Errorf("set header style %s: %w", cell, err)
		}
	}
	if err := f.SetRowHeight(sheet, 2, 40); err != nil {
		return fmt.Errorf("set header row height: %w", err)
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
		row := i + 3

		viral := ""
		if v.Viral {
			viral = "+"
		}
		relevant := "-"
		if v.IsRelevant {
			relevant = "+"
		}
		status := v.Status
		if status == "created" {
			status = ""
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
			status,
			deref(v.ErrorStage),
			deref(v.ErrorMessage),
			deref(v.AnalysisProvider),
			deref(v.RawPayload),
			deref(v.LLMProvider),
			deref(v.LLMModel),
			deref(v.Brief),
			deref(v.Script),
			deref(v.GenerationPlatform),
			deref(v.GenerationExternalID),
			deref(v.GenerationStatus),
			deref(v.GenerationErrorMessage),
			deref(v.GenerationS3URL),
		}

		for col, val := range vals {
			if err := set(row, col+1, val); err != nil {
				return err
			}
		}
	}

	for i := 3; i <= len(videos)+2; i++ {
		if err := f.SetRowHeight(sheet, i, 18); err != nil {
			return err
		}
	}

	return nil
}
