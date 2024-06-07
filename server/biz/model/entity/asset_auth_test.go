package entity

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSliceReflect(t *testing.T) {
	//var data []entity.SysUser
	var data = new([]SysUser)
	//var data []entity.SysUser
	valueOf := reflect.ValueOf(data)
	fmt.Printf("valueOf: %v\n", valueOf)
	valueOfType := valueOf.Type()
	fmt.Printf("valueOfType: %v\n", valueOfType)
	valueOfTypeElem := valueOfType.Elem()
	fmt.Printf("valueOfTypeElem: %v\n", valueOfTypeElem)
	valueOfTypeElemKind := valueOfTypeElem.Kind()
	fmt.Printf("valueOfTypeElemKind: %v\n", valueOfTypeElemKind)
	valueOfTypeElemElem := valueOfTypeElem.Elem()
	fmt.Printf("valueOfTypeElemElem: %v\n", valueOfTypeElemElem)
	valueOfTypeKind := valueOfType.Kind()
	fmt.Printf("valueOfTypeKind: %v\n", valueOfTypeKind)
	valueOfTypeKindString := valueOfTypeKind.String()
	fmt.Printf("valueOfTypeKindString: %v\n", valueOfTypeKindString)
	println("========================")
	valueOf = reflect.Indirect(valueOf)
	fmt.Printf("valueOf: %v\n", valueOf)
	valueOfType = valueOf.Type()
	fmt.Printf("valueOfType: %v\n", valueOfType)
	valueOfTypeElem = valueOfType.Elem()
	fmt.Printf("valueOfTypeElem: %v\n", valueOfTypeElem)
	valueOfTypeElemKind = valueOfTypeElem.Kind()
	fmt.Printf("valueOfTypeElemKind: %v\n", valueOfTypeElemKind)
	valueOfTypeKind = valueOfType.Kind()
	fmt.Printf("valueOfTypeKind: %v\n", valueOfTypeKind)
	valueOfTypeKindString = valueOfTypeKind.String()
	fmt.Printf("valueOfTypeKindString: %v\n", valueOfTypeKindString)
	println("========================")
	fmt.Printf("isSlice %v\n", reflect.Slice == reflect.Indirect(valueOf).Type().Kind())
	fmt.Printf("isSlice %v\n", reflect.Slice == reflect.Indirect(valueOf).Kind())

}

func TestStructReflect(t *testing.T) {
	var data SysUser
	//var data []entity.SysUser
	valueOf := reflect.ValueOf(data)
	fmt.Printf("valueOf: %v\n", valueOf)
	fmt.Printf("valueOfKind: %v\n", valueOf.Kind())
	valueOfType := valueOf.Type()
	fmt.Printf("valueOfType: %v\n", valueOfType)
	valueOfTypeKind := valueOfType.Kind()
	fmt.Printf("valueOfTypeKind: %v\n", valueOfTypeKind)
	valueOfTypeKindString := valueOfTypeKind.String()
	fmt.Printf("valueOfTypeKindString: %v\n", valueOfTypeKindString)
	println("========================")
	valueOf = reflect.Indirect(valueOf)
	fmt.Printf("valueOf: %v\n", valueOf)
	fmt.Printf("valueOfKind: %v\n", valueOf.Kind())
	valueOfType = valueOf.Type()
	fmt.Printf("valueOfType: %v\n", valueOfType)
	valueOfTypeKind = valueOfType.Kind()
	fmt.Printf("valueOfTypeKind: %v\n", valueOfTypeKind)
	valueOfTypeKindString = valueOfTypeKind.String()
	fmt.Printf("valueOfTypeKindString: %v\n", valueOfTypeKindString)
	println("========================")
	fmt.Printf("isSlice %v\n", reflect.Slice == reflect.Indirect(valueOf).Type().Kind())

}
