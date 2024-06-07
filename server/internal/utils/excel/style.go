package excel

import "github.com/xuri/excelize/v2"

func defaultBorder() []excelize.Border {
	return []excelize.Border{
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}
}

func defaultRowStyle() *excelize.Style {
	style := excelize.Style{}
	style.Border = defaultBorder()
	style.Alignment = &excelize.Alignment{
		// 水平对齐居中
		Horizontal: "center",
		// 垂直对齐居中
		Vertical: "center",
		// 自动换行
		WrapText: true,
	}
	style.Font = &excelize.Font{
		Size: 12,
	}
	return &style
}
