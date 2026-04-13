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
		sheet := strings.Title(platform)

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

func writeVideosSheet(f *excelize.File, sheet string, videos []*view.Video) error {
	headers := []string{
		"Title",
		"Views",
		"Likes",
		"Comments",
		"PublishedAt",
		"CreatedAt",
		"URL",
		"BloggerURL",
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

	for i, v := range videos {
		row := i + 2

		if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", row), v.Title); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Views); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Likes); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.Comments); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("E%d", row), v.PublishedAt.Format("2006-01-02 15:04")); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("F%d", row), v.CreatedAt.Format("2006-01-02 15:04")); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("G%d", row), v.URL); err != nil {
			return err
		}
	}

	return nil
}

