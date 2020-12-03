package sign

// 校验防止篡改和重放攻击

import (
	"fmt"
	"sort"
	"strings"

	"github.com/haobird/golibs/md5f"
)

//Encrypt 加密字符串
func Encrypt(params map[string]interface{}) string {
	// 摘取出AppSecret参数
	secretName := "AppSecret"
	return EncryptSpecial(params, secretName)
}

//EncryptSpecial 特殊指定
func EncryptSpecial(params map[string]interface{}, secretName string) string {
	secretStr := ""
	if val, exist := params[secretName]; exist {
		secretStr = "&" + secretName + "=" + val.(string)
		delete(params, secretName)
	}
	str := order(params)
	str = str + secretStr
	str = md5f.MD5(str)
	return strings.ToUpper(str)
}

//order
func order(params map[string]interface{}) string {
	var key []string
	var str = ""
	for k := range params {
		key = append(key, k)
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		if i == 0 {
			str = fmt.Sprintf("%v=%v", key[i], params[key[i]])
		} else {
			str = str + fmt.Sprintf("&%v=%v", key[i], params[key[i]])
		}
	}
	return str
}
