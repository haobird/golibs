package main

import (
	"fmt"

	"github.com/haobird/golibs/jsonf"
	"github.com/haobird/golibs/md5f"
)

// 测试
func main() {
	str := "test"
	md5 := md5f.MD5(str)
	fmt.Println(md5)

	var f map[string]interface{}
	text := `{
		"Nonce": "iVJZykib4H",
		"Status": 1,
		"AppKey": "a1668901ccf1a3f792361fb88d707ddd",
		"DeviceCode": "05CA19070001",
		"Time": "1606898562000",
		"Sign": "BF10EE4E0758A5A97DFE9DDE83F39F6E",
		"Timestamp": 1606900347844,
		"SpaceNo": "101001"
	  }`
	jsonf.DecodeUseNumber([]byte(text), &f)
	fmt.Println(f)
}
