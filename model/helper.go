package model

import (
	"crypto/md5"
	"fmt"
)



func EncodePassword(pwd string) string {
	h:=md5.New()
	h.Write([]byte(pwd+pwd+pwd))
	return fmt.Sprintf("%x",h.Sum(nil))
}