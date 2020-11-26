package utils

// 辅助函数
// 不依赖任何项目的独立辅助函数

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

//GenShortID 获取随机id
func GenShortID() (string, error) {
	return shortid.Generate()
}

//GetReqID 获取requestid
func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestID, ok := v.(string); ok {
		return requestID
	}
	return ""
}

//Trim 清除字符串两边空格
func Trim(str string) string {
	strList := []byte(str)
	spaceCount, count := 0, len(strList)
	for i := 0; i <= len(strList)-1; i++ {
		if strList[i] == 32 {
			spaceCount++
		} else {
			break
		}
	}

	strList = strList[spaceCount:]
	spaceCount, count = 0, len(strList)
	for i := count - 1; i >= 0; i-- {
		if strList[i] == 32 {
			spaceCount++
		} else {
			break
		}
	}

	return string(strList[:count-spaceCount])
}

//Request 发起 Http请求
func Request(url string, method string, data []byte, headers map[string]string) (result string, err error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodystr := string(body)
	return bodystr, nil

	// if resp.StatusCode == 200 {
	// 	body, err_ := ioutil.ReadAll(resp.Body)
	// 	if err_ != nil {
	// 		return "", err_
	// 	}
	// 	bodystr := string(body)
	// 	return bodystr, nil
	// }
	// return "", err

}

//ParseRequestPacket 处理http请求
func ParseRequestPacket(buf []byte) (*http.Request, error) {
	ioReader := bytes.NewReader(buf)
	reader := bufio.NewReader(ioReader)
	req, err := http.ReadRequest(reader)
	// fmt.Println(req.URL)
	return req, err
}

/** * 字符串首字母转化为大写 ios_bbbbbbbb -> iosBbbbbbbbb */
func strFirstToUpper(str string) string {
	temp := strings.Split(str, ".")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				vv[i] -= 32
				upperStr += string(vv[i]) // + string(vv[i+1])
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	return upperStr
}

//dont do this, see above edit
func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
