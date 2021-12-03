package tools

import (
	"bytes"
	"math/big"
	"reflect"
	"sort"
)

func SortStrings(strs []string)  {
	sort.Strings(strs)
}

func BytesCmp(bs1,bs2 []byte)bool{
	b1 := big.NewInt(0).SetBytes(bs1)
	b2 := big.NewInt(0).SetBytes(bs2)
	if b1.Cmp(b2)>0{
		return true
	}
	return false
}

func BytesEqual(bs1,bs2 []byte) bool{
	b1 := big.NewInt(0).SetBytes(bs1)
	b2 := big.NewInt(0).SetBytes(bs2)
	if b1.Cmp(b2)==0{
		return true
	}
	return false
}

//ratio -1 小于 0 相等 1大于
func valueRatio(iValue,jValue reflect.Value)(ratio,equal bool) {
	switch iValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iv := iValue.Int()
		jv := jValue.Int()
		return iv < jv,iv == jv
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		iv := iValue.Uint()
		jv := jValue.Uint()
		return iv < jv,iv == jv
	case reflect.Float32,reflect.Float64:
		iv := iValue.Float()
		jv := jValue.Float()
		return iv < jv,iv == jv
	case reflect.String:
		iv := iValue.String()
		jv := jValue.String()
		return iv < jv,iv == jv
	case reflect.Slice:
		iv := iValue.Bytes()
		jv := jValue.Bytes()
		return BytesCmp(iv,jv),bytes.Equal(iv,jv)
	default:
		panic("unknown type")
	}
}

func SortStructList(list interface{}, tags []string){
	sort.Slice(list, func(i, j int) bool {
		for _,tag := range tags{
			var iValue,jValue reflect.Value
			if reflect.ValueOf(list).Index(i).Kind() == reflect.Ptr{
				iValue = reflect.ValueOf(list).Index(i).Elem().FieldByName(tag)
			}else {
				iValue = reflect.ValueOf(list).Index(i).FieldByName(tag)
			}
			if reflect.ValueOf(list).Index(j).Kind() == reflect.Ptr{
				jValue = reflect.ValueOf(list).Index(j).Elem().FieldByName(tag)
			}else {
				jValue = reflect.ValueOf(list).Index(j).FieldByName(tag)
			}
			//iValue := reflect.ValueOf(list).Index(i).FieldByName(tag)
			//jValue := reflect.ValueOf(list).Index(j).FieldByName(tag)
			if iValue.Kind() != jValue.Kind() {
				panic("SortList Inconsistent type")
			}
			b,equal := valueRatio(iValue,jValue)
			if !equal{
				return b
			}
		}
		return false
	})
	return
}