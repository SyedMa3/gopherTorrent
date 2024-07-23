package bencode

import "fmt"

func encodeBencode(s BencodeValue) []byte {
	switch v := s.(type) {
	case int64:
		return encodeInt(s)
	case string:
		return encodeString(s)
	case map[string]BencodeValue:
		return encodeDict(s.(map[string]BencodeValue))
	case []BencodeValue:
		return encodeList(s.([]BencodeValue))
	default:
		fmt.Errorf("unkown decoded type %s", v)
		return nil
	}
}

func encodeInt(s BencodeValue) []byte {
	var num string

	num += "i"
	num += string(s.(int64))
	num += "e"

	return []byte(num)
}

func encodeString(s BencodeValue) []byte {
	t := s.(string)
	var str string
	str += string(len(t)) + ":"
	str += t

	return []byte(str)
}

func encodeList(l []BencodeValue) []byte {
	var str string
	str += "l"
	for _, i := range l {
		str += string(encodeBencode(i))
	}
	str += "e"

	return []byte(str)
}

func encodeDict(d map[string]BencodeValue) []byte {
	var str string
	str += "d"
	for k, v := range d {
		str += string(encodeBencode(k))
		str += string(encodeBencode(v))
	}
	str += "e"

	return []byte(str)
}
