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
		"BloggerURL",
		"Title",
		"URL",
		"PublishedAt",
		"CreatedAt",
		"Views",
		"Likes",
		"Comments",
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

	for i, v := range videos {
		row := i + 2

		if err := set(row, 1, v.BloggerURL); err != nil {
			return err
		}
		if err := set(row, 2, v.Title); err != nil {
			return err
		}
		if err := set(row, 3, v.URL); err != nil {
			return err
		}
		if err := set(row, 4, v.PublishedAt.Format("2006-01-02")); err != nil {
			return err
		}
		if err := set(row, 5, v.CreatedAt.Format("2006-01-02")); err != nil {
			return err
		}
		if err := set(row, 6, v.Views); err != nil {
			return err
		}
		if err := set(row, 7, v.Likes); err != nil {
			return err
		}
		if err := set(row, 8, v.Comments); err != nil {
			return err
		}
	}

	return nil
}

