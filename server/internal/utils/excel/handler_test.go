package excel

import (
	"fmt"
	"go-protector/server/internal/custom/c_type"
	"os"
	"reflect"
	"testing"
)

type TestRowStruct struct {
	//LineNum int
	//ErrMsg  string `excel:"title:错误消息;width:50"`
	StdRow
	Name  string      `excel:"title:名称"`
	Age   int         `excel:"title:年龄"`
	Email c_type.Time `excel:"title:邮箱" binding:"required"`
}

//
//func (_self *TestRowStruct) GetLineNum() int {
//	return _self.LineNum
//}
//
//func (_self *TestRowStruct) SetLineNum(i int) {
//	_self.LineNum = i
//}
//
//func (_self *TestRowStruct) SetErr(err error) {
//	_self.ErrMsg = err.Error()
//}
//
//func (_self *TestRowStruct) GetErr() string {
//	return _self.ErrMsg
//}

type TestHandler[T TestRowStruct] struct {
	ErrData []TestRowStruct
}

func (_self *TestHandler[T]) ReadRow(row *TestRowStruct) error {
	fmt.Printf("ReadRow: %v\n", *row)
	return nil
}

func (_self *TestHandler[T]) ReadDone() {
	fmt.Printf("ReadDone")

}

func (_self *TestHandler[T]) NewRow() *TestRowStruct {
	return &TestRowStruct{}
}

func (_self *TestHandler[T]) AppendErrData(row *TestRowStruct) {
	_self.ErrData = append(_self.ErrData, *row)
}

func Test_analysisRowStruct(t *testing.T) {
	type args struct {
		rowStruct interface{}
	}
	tests := []struct {
		name          string
		args          args
		wantFileSlice []ColInfo
	}{
		{
			name: "test",
			args: args{rowStruct: TestRowStruct{}},
			wantFileSlice: []ColInfo{
				{
					Name:  "Name",
					Title: "名称",
					Index: 0,
				}, {
					Name:  "Age",
					Title: "年龄",
					Index: 1,
				}, {
					Name:  "Email",
					Title: "邮箱",
					Index: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFileSlice := analysisRowStruct(tt.args.rowStruct); !reflect.DeepEqual(gotFileSlice, tt.wantFileSlice) {
				t.Errorf("analysisRowStruct() = %v, want %v", gotFileSlice, tt.wantFileSlice)
			}
		})
	}
}

func TestGenerateExcel(t *testing.T) {
	var slice []TestRowStruct
	slice = append(slice, TestRowStruct{
		Name: "name1",
		Age:  1,
		//Email: "email1",
	}, TestRowStruct{
		Name: "name2",
		Age:  2,
		//Email: "email2",
	})

	file, err := GenerateExcel(slice, "错误消息")
	if err != nil {
		t.Errorf("err: %v\n", err)
		return
	}
	f, _ := os.Create("/opt/work_space/github/go-protector/data/test.xlsx")
	defer f.Close()
	if err = file.Write(f); err != nil {
		t.Errorf("err: %v\n", err)

		return
	}

}

func TestReadExcel(t *testing.T) {

	f, err := os.Open("/opt/work_space/github/go-protector/data/test.xlsx")
	if err != nil {
		t.Errorf("open err: %v", err)
		return
	}
	handler := TestHandler[TestRowStruct]{}
	err = ReadExcelFirstSheet[*TestRowStruct](f, &handler)
	if err != nil {
		t.Errorf("ReadExcelFirstSheet err: %v", err)
		return
	}
	if len(handler.ErrData) <= 0 {
		return
	}

	file, err := GenerateExcel(handler.ErrData)

	if err != nil {
		t.Errorf("err: %v\n", err)
		return
	}
	defer file.Close()
	f, _ = os.Create("/opt/work_space/github/go-protector/data/test_err.xlsx")
	defer f.Close()
	if err = file.Write(f); err != nil {
		t.Errorf("err: %v\n", err)
		return
	}

}
