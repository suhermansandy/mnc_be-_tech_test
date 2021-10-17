package excel

import (
	"io"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ExcelHandler interface {
	NewFile() ExcelHandler
	SetSheetName(oldName string, newName string)
	GetSheetName(index int) string
	SetCellValue(sheet string, axis string, value interface{})
	AutoFilter(sheet string, hcell string, vcell string, format string) error
	Write(w io.Writer) error
}

type xlsx struct {
	excelizeFile *excelize.File
}

func New() ExcelHandler {
	return &xlsx{}
}

func Wrap(file *excelize.File) ExcelHandler {
	return &xlsx{file}
}

func (x *xlsx) NewFile() ExcelHandler {
	file := excelize.NewFile()
	return Wrap(file)
}

func (x *xlsx) SetSheetName(oldName string, newName string) {
	x.excelizeFile.SetSheetName(oldName, newName)
}

func (x *xlsx) GetSheetName(index int) string {
	return x.excelizeFile.GetSheetName(index)
}

func (x *xlsx) SetCellValue(sheet string, axis string, value interface{}) {
	x.excelizeFile.SetCellValue(sheet, axis, value)
}

func (x *xlsx) AutoFilter(sheet string, hcell string, vcell string, format string) error {
	return x.excelizeFile.AutoFilter(sheet, hcell, vcell, format)
}

func (x *xlsx) Write(w io.Writer) error {
	return x.excelizeFile.Write(w)
}
