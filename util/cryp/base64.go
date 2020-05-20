package cryp

import (
	b64 "encoding/base64"
)

func Base64Encode(s string) string {
	return b64.StdEncoding.EncodeToString([]byte(s))
}
func Base64EncodeByte(bs []byte) string {
	return b64.StdEncoding.EncodeToString(bs)
}
func Base64Decode(s string) (string, error) {
	ds, err := b64.StdEncoding.DecodeString(s)
	return string(ds), err
}
func Base64DecodeByte(s string) []byte {
	ds, _ := b64.StdEncoding.DecodeString(s)
	return ds
}
func Base64EncodeByteCount(bs []byte, count int) string {
	s := Base64EncodeByte(bs)
	if count < 2 {
		return s
	}
	for i := 0; i < count-1; i++ {
		s = Base64Encode(s)
	}
	return s
}
