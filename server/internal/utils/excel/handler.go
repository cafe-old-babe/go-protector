package excel

import (
	"cmp"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"go-protector/server/internal/custom/c_translator"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/utils"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var tagExcel = "excel"

var title = "title"
var index = "index"
var width = "width"

var typeOfTime = reflect.TypeOf(c_type.Time{})

type Handler[T Row] interface {
	// ReadRow 处理每行数据
	ReadRow(row T) error
	// ReadDone 最后回调
	ReadDone()
	// NewRow 获取结构体
	NewRow() T
	// AppendErrData 追加错误数据
	AppendErrData(row T)
}

// ReadExcelFirstSheet 读取文件  https://blog.csdn.net/qq_39272466/article/details/131631964
func ReadExcelFirstSheet[T Row](readCloser io.ReadCloser, h Handler[T]) (err error) {
	defer readCloser.Close()

	file, err := excelize.OpenReader(readCloser)
	if err != nil {
		return
	}
	defer file.Close()
	sheetName := file.GetSheetName(0)
	if sheetName == "" {
		err = errors.New("没有sheet")
		return
	}
	rows, err := file.Rows(sheetName)
	if err != nil {
		return
	}

	colInfos := analysisRowStruct(h.NewRow())

	var titleCols []string

	if !rows.Next() {
		err = errors.New("请下载模板再导入")
		return
	}
	if titleCols, err = rows.Columns(); err != nil {
		return
	}
	if len(titleCols) != len(colInfos) {
		err = errors.New("请下载模板再导入")
		return

	}
	for i, tc := range titleCols {
		// 校验
		if tc != colInfos[i].Title {
			err = errors.New("请下载模板再导入")
			return
		}
	}
	processRows[T](rows, colInfos, h)

	return
}

// processRows 处理每行数据
func processRows[T Row](rows *excelize.Rows, colInfos []ColInfo, h Handler[T]) {

	defer rows.Close()
	defer h.ReadDone()

	var columns []string
	var err error
	var lineNum int
	var colInfo ColInfo
	var rowData T
	var indirect reflect.Value

	for rows.Next() {
		func() {
			err = nil
			defer func() {
				defer func() {
					if err != nil {
						err = c_translator.ConvertValidateErr(err)
						rowData.SetErr(err)
						h.AppendErrData(rowData)
						return
					}
				}()
				if recoverErr := recover(); recoverErr != nil {
					err = errors.New(recoverErr.(string))
					return
				}

				err = h.ReadRow(rowData)
			}()
			lineNum++
			rowData = h.NewRow()
			rowData.SetLineNum(lineNum)

			if columns, err = rows.Columns(); err != nil {
				return
			}
			// 校验必填项
			for i, column := range columns {
				column = strings.Trim(column, " ")
				colInfo = colInfos[i]
				if len(column) <= 0 {
					continue
				}

				indirect = reflect.Indirect(reflect.ValueOf(rowData))
				fieldStruct, _ := indirect.Type().FieldByName(colInfo.Name)

				// 将field的值设置为 column
				fieldVal := indirect.FieldByName(colInfo.Name)
				if !fieldVal.CanSet() {
					continue
				}
				switch fieldStruct.Type.Kind() {
				case reflect.String:
					fieldVal.SetString(column)
				case reflect.Int:
					if reflect.ValueOf(fieldVal.Interface()).CanInt() {
						var v int64
						if v, err = strconv.ParseInt(column, 10, 64); err != nil {
							err = errors.Join(errors.New(colInfo.Title + "赋值失败,请联系管理员"))
						} else {
							fieldVal.SetInt(v)
						}
					}
				case reflect.Struct:
					field := fieldStruct.Type
					if field.AssignableTo(typeOfTime) {
						if time, parseErr := utils.ParseTime(column); parseErr == nil {
							fieldVal.Set(reflect.ValueOf(time))
						} else {
							err = errors.Join(errors.New(colInfo.Title + ":[" + column + "].时间格式不正确:" + parseErr.Error()))
						}
					} else {
						err = errors.Join(errors.New(colInfo.Title + "赋值失败,请联系管理员"))
					}

				default:
					err = errors.Join(errors.New(colInfo.Title + "赋值失败,请联系管理员"))
				}

			}
		}()
	}

}

// analysisRowStruct 解析结构体
func analysisRowStruct(rowStruct interface{}) (colInfos []ColInfo) {

	indirect := reflect.Indirect(reflect.ValueOf(rowStruct))
	indirectType := indirect.Type()

	var elemSlice []string
	var elemKV []string
	var elemV string
	var err error
	var atoi int
	for i := 0; i < indirectType.NumField(); i++ {
		field := indirectType.Field(i)
		if !field.IsExported() {
			continue
		}
		if field.Anonymous {

			colInfos = append(colInfos, analysisRowStruct(indirect.Field(i).Interface())...)
			continue
		}
		tagVal, ok := field.Tag.Lookup(tagExcel)
		if !ok {
			continue
		}
		if tagVal = strings.Trim(tagVal, " "); len(tagVal) <= 0 {
			continue
		}
		if elemSlice = strings.Split(tagVal, ";"); len(elemSlice) <= 0 {
			continue
		}

		elemMap := map[string]string{}

		for _, val := range elemSlice {
			if len(val) <= 0 {
				continue
			}
			if elemKV = strings.Split(val, ":"); len(elemKV) <= 0 || len(elemKV) > 2 {
				continue
			}
			if len(elemKV) == 1 {
				elemMap[elemKV[0]] = ""
			} else {
				elemMap[elemKV[0]] = elemKV[1]
			}
		}
		if len(elemMap) <= 0 {
			continue
		}
		colInfo := ColInfo{}
		colInfo.Name = field.Name
		elemV = elemMap[title]
		if elemV = strings.Trim(elemV, " "); len(elemV) <= 0 {
			elemV = field.Name
		}
		colInfo.Title = elemV

		elemV = elemMap[index]
		if elemV = strings.Trim(elemV, " "); len(elemV) <= 0 {
			atoi = i
		} else {
			if atoi, err = strconv.Atoi(elemV); err != nil {
				atoi = i
			}
		}
		colInfo.Index = atoi

		elemV = elemMap[width]
		if elemV = strings.Trim(elemV, " "); len(elemV) <= 0 {
			atoi = 25
		} else {
			if atoi, err = strconv.Atoi(elemV); err != nil {
				atoi = 25
			}
		}
		colInfo.Width = atoi

		colInfos = append(colInfos, colInfo)

	}
	slices.SortFunc(colInfos, func(a, b ColInfo) int {
		return cmp.Compare(a.Index, b.Index)
	})

	return
}

// GenerateExcel 生成excel文件
func GenerateExcel(sliceOrStruct interface{}, hideTitle ...string) (file *excelize.File, err error) {
	defer func() {
		if err != nil && file != nil {
			_ = file.Close()
			file = nil
		}
	}()
	indirect := reflect.Indirect(reflect.ValueOf(sliceOrStruct))
	var hasData bool
	var value reflect.Value

	switch indirect.Kind() {
	case reflect.Slice:
		hasData = indirect.Len() > 0
		value = reflect.New(indirect.Type().Elem()).Elem()
	case reflect.Struct:
		hasData = !indirect.IsZero()
		value = reflect.New(indirect.Type()).Elem()

	default:
		err = errors.New("数据类型不正确")
		return
	}

	colStructs := analysisRowStruct(value.Interface())
	if len(colStructs) <= 0 {
		err = errors.New("无法匹配title")
		return
	}
	var data [][]string
	var elemValue reflect.Value
	var fieldValue reflect.Value
	var colVal string
	data = append(data, utils.SliceToFieldSlice[string]("Title", colStructs))

	if !hasData {
		goto dataLabel
	}

	for i := 0; i < indirect.Len(); i++ {
		elemValue = indirect.Index(i)
		if elemValue.IsZero() {
			continue
		}
		var row []string
		for _, rs := range colStructs {
			fieldValue = elemValue.FieldByName(rs.Name)
			if fieldValue.IsZero() {
				colVal = ""
			} else {
				if fieldValue.Type().AssignableTo(typeOfTime) {
					time := fieldValue.Interface().(c_type.Time)
					if !time.Valid {
						colVal = ""
					} else {
						colVal = time.String()
					}
				} else {

					colVal = fmt.Sprintf("%v", fieldValue)
				}

			}
			row = append(row, colVal)
		}
		if len(row) > 0 {
			data = append(data, row)
		}
	}
dataLabel:
	if len(data) <= 0 {
		err = errors.New("无数据")
		return
	}
	file = excelize.NewFile()
	var sheetNum int
	sheetName := "Sheet1"
	if sheetNum, err = file.NewSheet(sheetName); err != nil {
		return
	}
	var rowStyle int
	if rowStyle, err = file.NewStyle(defaultRowStyle()); err != nil {
		return
	}
	var topLeftCell, bottomRightCell string
	if topLeftCell, err = excelize.CoordinatesToCellName(1, 1); err != nil {
		return
	}
	if bottomRightCell, err = excelize.CoordinatesToCellName(len(data[0]), len(data)); err != nil {
		return
	}
	if err = file.SetCellStyle(sheetName, topLeftCell, bottomRightCell, rowStyle); err != nil {
		return
	}
	var cellName string
	for i, colStruct := range colStructs {
		if cellName, err = excelize.ColumnNumberToName(i + 1); err != nil {
			return
		}
		if err = file.SetColWidth(sheetName, cellName, cellName, float64(colStruct.Width)); err != nil {
			return
		}
		if slices.Contains(hideTitle, colStruct.Title) {
			if err = file.SetColVisible(sheetName, cellName, false); err != nil {
				return nil, err
			}
		}
	}
	for i, row := range data {
		if cellName, err = excelize.CoordinatesToCellName(1, i+1); err != nil {
			return
		}

		if err = file.SetSheetRow(sheetName, cellName, &row); err != nil {
			return
		}

	}

	file.SetActiveSheet(sheetNum)
	return
}
