package signer

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// 签名
// 排序
// 凡是传递过来的，都要加密，除了sign和appsecret字段

//Signer 签名
type Signer struct {
	DefaultKeyName
	// bodyPrefix    string     // 参数体前缀
	// bodySuffix    string     // 参数体后缀
	secretVal     string        // 签名密钥
	body          url.Values    // 签名参数体
	otherFields   []string      // 其它需要签名的字段
	skippedFields []string      // 需要忽略的签名字段
	timeout       time.Duration // 签名过期时间

}

//Init 初始化相关字段
func (slf *Signer) Init(values url.Values) {
	slf.ParseValues(values)
	// 完善需要加签的字段
	if _, ok := slf.body[slf.keyNameTimestamp]; !ok {
		slf.SetTimeStamp(time.Now().Unix())
	}

	if _, ok := slf.body[slf.keyNameNonceStr]; !ok {
		slf.RandNonceStr()
	}

	if slf.secretVal != "" {
		slf.AddBody(slf.keyNameAppSecret, slf.secretVal)
	} else if val, ok := slf.body[slf.keyNameAppSecret]; ok {
		slf.secretVal = val[0]
	}
}

// ParseValues 将Values参数列表解析成参数Map。如果参数是多值的，则将它们以逗号Join成字符串。
func (slf *Signer) ParseValues(values url.Values) {
	for key, value := range values {
		slf.body[key] = value
	}
}

// AddBody 添加签名体字段和值
func (slf *Signer) AddBody(key string, value string) *Signer {
	return slf.AddBodies(key, []string{value})
}

// AddBodies 添加签名字段
func (slf *Signer) AddBodies(key string, value []string) *Signer {
	slf.body[key] = value
	return slf
}

// SetAppSecret 设置签名密钥
func (slf *Signer) SetAppSecret(appSecret string) *Signer {
	slf.secretVal = appSecret
	return slf
}

// GetAppSecret 返回AppSecret内容
func (slf *Signer) GetAppSecret() string {
	// 如果存在，则设置
	return slf.body.Get(slf.keyNameAppSecret)
}

// SetTimeStamp 设置时间戳参数
func (slf *Signer) SetTimeStamp(ts int64) *Signer {
	return slf.AddBody(slf.keyNameTimestamp, strconv.FormatInt(ts, 10))
}

// GetTimestamp 获取TimeStamp
func (slf *Signer) GetTimestamp() string {
	return slf.body.Get(slf.keyNameTimestamp)
}

// GetSign 获取签名
func (slf *Signer) GetSign() string {
	return slf.body.Get(slf.keyNameSign)
}

// SetNonceStr 设置随机字符串参数
func (slf *Signer) SetNonceStr(nonce string) *Signer {
	return slf.AddBody(slf.keyNameNonceStr, nonce)
}

// RandNonceStr 自动生成16位随机字符串参数
func (slf *Signer) RandNonceStr() *Signer {
	return slf.SetNonceStr(RandString(16))
}

// GetBodyString 获取用于签名的原始字符串
func (slf *Signer) GetBodyString() string {
	return slf.getSortedBodyString()
}

// GetSignRawString 获取未签名前的字符串
func (slf *Signer) GetSignRawString() string {
	bodyStr := slf.GetBodyString()
	// secret := slf.GetAppSecret()
	secret := slf.secretVal
	fmt.Println("密钥")
	fmt.Println(secret)
	if secret != "" {
		bodyStr = bodyStr + "&" + slf.keyNameAppSecret + "=" + secret
	}
	return bodyStr
}

// GetSignedQuery 获取带签名参数的字符串
func (slf *Signer) GetSignedQuery() string {
	body := slf.GetBodyString()
	sign := slf.GetSignature()
	return body + "&" + slf.keyNameSign + "=" + sign
}

// GetSignature 获取签名
func (slf *Signer) GetSignature() string {
	sign := encrypt(slf.GetSignRawString())
	return sign
}

// getSortedBodyString 获取排序好的字符串
func (slf *Signer) getSortedBodyString() string {
	// 删除需要忽略的字段
	validBody := slf.FilterFields()
	// 校验必含的字段是否完整
	return SortKVPairs(validBody)
}

// Debug 调试
func (slf *Signer) Debug() map[string]string {
	var arr map[string]string
	arr["sortBodyStr"] = slf.GetBodyString()
	arr["signature"] = slf.GetSignature()
	arr["signRawString"] = slf.GetSignRawString()
	arr["signedQuery"] = slf.GetSignedQuery()

	return arr
}

// GetBody 返回Body内容
func (slf *Signer) GetBody() url.Values {
	return slf.body
}

//FilterFields 过滤字段
func (slf *Signer) FilterFields() url.Values {
	var validBody map[string][]string
	// 如果指定了加签名字段，则优先使用加签名字段
	if len(slf.otherFields) > 0 {
		fields := slf.MustHasKeys()
		for _, key := range fields {
			if val, hit := slf.body[key]; hit {
				validBody[key] = val
			}
		}
	} else if len(slf.skippedFields) > 0 {
		validBody = slf.body
		for _, key := range slf.skippedFields {
			delete(validBody, key)
		}
	} else {
		validBody = slf.body
	}
	// 过滤掉sign字段
	delete(validBody, slf.keyNameSign)
	delete(validBody, slf.keyNameAppSecret)

	// 如果没有指定加签字段，则执行过滤字段之后加签
	return validBody
}

// CheckTimeStamp 检查时间戳有效期
// 10位数的时间戳是以 秒 为单位
// 13位数的时间戳是以 毫秒 为单位
// 19位数的时间戳是以 纳秒 为单位
func (slf *Signer) CheckTimeStamp() error {
	timestamp := slf.GetTimestamp()
	length := len(timestamp)
	// 格式化为int64
	timeInt, _ := strconv.ParseInt(timestamp, 10, 64)
	thatTime := time.Unix(timeInt, 0)
	if length > 10 {
		thatTime = time.Unix(0, timeInt*int64(time.Millisecond))
	}

	fmt.Println(thatTime)
	if time.Now().Sub(thatTime) > slf.timeout {
		return fmt.Errorf("TIMESTAMP_TIMEOUT:<%s>", timestamp)
	}
	return nil
}

// MustHasKeys 必须包含指定的字段参数
func (slf *Signer) MustHasKeys() []string {
	fields := []string{slf.keyNameTimestamp, slf.keyNameNonceStr, slf.keyNameSign, slf.keyNameAppID}
	fields = append(slf.otherFields, fields...)
	return fields
}

//CheckMustHasKeys 校验必选字段
func (slf *Signer) CheckMustHasKeys() error {
	fields := slf.MustHasKeys()
	for _, key := range fields {
		if _, hit := slf.body[key]; !hit {
			return fmt.Errorf("KEY_MISSED:<%s>", key)
		}
	}
	return nil
}

//CheckSignature 校验签名
func (slf *Signer) CheckSignature() error {
	// 获取 当前的签名
	oldSign := slf.GetSign()
	// 生成新的签名
	newSign := slf.GetSignature()
	// 比较判断
	if oldSign == "" || oldSign != newSign {
		return fmt.Errorf("sign Error:<%s>", oldSign)
	}
	return nil
}
