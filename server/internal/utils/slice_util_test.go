package utils

import (
	"fmt"
	"go-protector/server/biz/model/entity"
	"reflect"
	"testing"
)

func TestSub(t *testing.T) {
	v1 := []uint64{1, 2, 3, 4, 9, 12, 32, 212, 412, 45231, 123, 31, 232, 123, 456, 2, 645, 2, 33, 124, 65, 23, 123, 623, 454, 22, 6246, 2, 4, 12524672, 451, 25, 26, 2, 521, 5, 6, 23, 51235}
	v2 := []uint64{1, 2, 3, 4, 9, 12, 32, 212, 412, 123, 31, 232, 123, 456, 2, 645, 2, 33, 124, 65, 23, 123, 623, 454, 22, 6246, 2, 4, 12524672, 451, 25, 26, 2, 521, 5, 6, 23, 51235}
	sub := SliceSub(v1, v2)
	fmt.Printf("SliceSub-->%v\n", sub)
	sub = SliceSubN(v1, v2)
	fmt.Printf("SliceSubN2-->%v\n", sub)
	del(&v1)
	m := map[uint64]uint64{1: 1}

	fmt.Printf("del-->%v\n", v1)
	delM(m)
	fmt.Printf("m 1 --> %v\n", m[1])

}

// BenchmarkSub  281751	      3957 ns/op
func BenchmarkSub(b *testing.B) {
	v1 := []uint64{1, 2, 3, 4, 9, 12, 32, 212, 412, 45231, 123, 31, 232, 123, 456, 2, 645, 2, 33, 124, 65, 23, 123, 623, 454, 22, 6246, 2, 4, 12524672, 451, 25, 26, 2, 521, 5, 6, 23, 51235}
	v2 := []uint64{1, 2, 3, 4, 9, 12, 32, 212, 412, 123, 31, 232, 123, 456, 2, 645, 2, 33, 124, 65, 23, 123, 623, 454, 22, 6246, 2, 4, 12524672, 451, 25, 26, 2, 521, 5, 6, 23, 51235}
	for i := 0; i < b.N; i++ {
		_ = SliceSub[uint64](v1, v2)
	}
}

// BenchmarkSubN2
func BenchmarkSubN2(b *testing.B) {
	v1 := []uint64{1, 2, 3, 4, 9, 12, 32, 212, 412, 45231, 123, 31, 232, 123, 456, 2, 645, 2, 33, 124, 65, 23, 123, 623, 454, 22, 6246, 2, 4, 12524672, 451, 25, 26, 2, 521, 5, 6, 23, 51235}
	v2 := []uint64{1, 2, 3, 4, 9, 12, 32, 212, 412, 123, 31, 232, 123, 456, 2, 645, 2, 33, 124, 65, 23, 123, 623, 454, 22, 6246, 2, 4, 12524672, 451, 25, 26, 2, 521, 5, 6, 23, 51235}
	for i := 0; i < b.N; i++ {
		_ = SliceSubN[uint64](v1, v2)
	}

}

func del(v *[]uint64) {
	// 删除 v 第一个元素
	*v = (*v)[1:]
}

func delM(m map[uint64]uint64) {
	delete(m, 1)
}

func TestSliceToFieldSlice(t *testing.T) {
	type args struct {
		slice []any
		field string
	}
	type testCase[T interface{ uint64 | string }] struct {
		name           string
		args           args
		wantFieldSlice []T
	}
	tests := []testCase[string]{

		{
			name: "1",
			args: args{
				slice: []any{
					entity.AssetAccount{
						AssetId: 0,
						Account: "12312",
					}},
				field: "Account",
			},
			wantFieldSlice: []string{"12312"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFieldSlice := SliceToFieldSlice[string](tt.args.field, tt.args.slice); !reflect.DeepEqual(gotFieldSlice, tt.wantFieldSlice) {
				t.Errorf("SliceToFieldSlice() = %v, want %v", gotFieldSlice, tt.wantFieldSlice)
			}
		})
	}
}
