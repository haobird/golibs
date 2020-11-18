package utils

// 辅助函数
// 不依赖任何项目的独立辅助函数

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

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

//ImageToBase64 图片转64
func ImageToBase64(path string, isUrl bool, imageType ...string) (string, error) {
	var (
		stream io.Reader
		err    error
	)
	if isUrl {
		var res *http.Response
		res, err = http.Get(path)
		if err == nil {
			stream = res.Body
		}
	} else {
		stream, err = os.Open(path)
	}
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(stream)
	imTyp := "jpeg"
	if imageType != nil {
		imTyp = imageType[0]
	}
	switch imTyp {
	case "jpeg":
		data, err = ToJpeg(data)
	case "png":
		data, err = ToPng(data)
	}
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

//GetMinioUrl 获取文件地址
func GetMinioUrl(fileName string) string {
	if fileName == "" {
		return ""
	}
	conf := viper.GetStringMap("minio")
	return fmt.Sprintf("http://%s/%s/%s", conf["endpoint"], conf["bucket"], fileName)
}

func ToPng(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/png":
		return imageBytes, nil
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("unable to convert %#v to png", contentType)
}

func ToJpeg(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/png":
		img, err := png.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case "image/jpeg":
		return imageBytes, nil
	}
	return nil, fmt.Errorf("unable to convert %#v to jpeg", contentType)
}
