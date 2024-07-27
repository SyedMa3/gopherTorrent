package bencode

import (
	"fmt"
	"reflect"
	"strconv"
)

func encodeInfo(info bencodeInfo) []byte {
	var enc []byte

	values := reflect.ValueOf(info)
	types := values.Type()
	enc = append(enc, "d"...)
	for i := 0; i < values.NumField(); i++ {
		// ty := types.Field(i).Type
		keyString := types.Field(i).Name
		if keyString == "pieceLength" {
			keyString = "piece length"
		}
		key := encodeString(reflect.ValueOf(keyString))
		value := encodeBencode(values.Field(i))
		enc = append(enc, key...)
		enc = append(enc, value...)
	}
	enc = append(enc, "e"...)

	return enc
}

func encodeBencode(s reflect.Value) []byte {
	switch s.Kind() {
	case reflect.Int64:
		return encodeInt(s)
	case reflect.String:
		return encodeString(s)
	case reflect.Map:
		return encodeDict(s.Interface().(map[string]BencodeValue))
	case reflect.Slice:
		return encodeList(s.Interface().([]BencodeValue))
	default:
		fmt.Println(reflect.TypeOf(s))
		fmt.Errorf("unkown decoded type %s", s)
		return nil
	}
}

func encodeInt(s reflect.Value) []byte {
	var num string

	num += "i"
	num += strconv.Itoa(int(s.Int()))
	num += "e"

	return []byte(num)
}

func encodeString(s reflect.Value) []byte {
	t := s.String()

	var str string
	str += strconv.Itoa(len(t)) + ":"
	str += t

	return []byte(str)
}

func encodeList(l []BencodeValue) []byte {
	var str string
	str += "l"
	for _, i := range l {
		str += string(encodeBencode(reflect.ValueOf(i)))
	}
	str += "e"

	return []byte(str)
}

func encodeDict(d map[string]BencodeValue) []byte {
	var str string
	str += "d"
	for k, v := range d {
		str += string(encodeBencode(reflect.ValueOf(k)))
		str += string(encodeBencode(reflect.ValueOf(v)))
	}
	str += "e"

	return []byte(str)
}
