package signer

import (
	"bytes"
	"encoding/json"
	"net/url"
	"time"
)

// 签名校验器

const (
	//KeyNameTimeStamp 时间戳关键字
	KeyNameTimeStamp = "timestamp"
	//KeyNameNonceStr 随机字符串
	KeyNameNonceStr = "nonce"
	//KeyNameAppID appid
	KeyNameAppID = "appid"
	//KeyNameAppSecret appsecret
	KeyNameAppSecret = "appsecret"
	//KeyNameSign sign
	KeyNameSign = "sign"
)

//DefaultKeyName 默认key的名字
type DefaultKeyName struct {
	keyNameTimestamp string
	keyNameNonceStr  string
	keyNameAppID     string
	keyNameSign      string
	keyNameAppSecret string
}

//Option 配置
type Option struct {
	KeyNameTimestamp string
	KeyNameNonceStr  string
	KeyNameAppID     string
	KeyNameSign      string
	KeyNameAppSecret string
	// bodyPrefix    string        // 参数体前缀
	// bodySuffix    string        // 参数体后缀
	SecretVal     string        // 签名密钥
	OtherFields   []string      // 其它需要签名的字段
	SkippedFields []string      // 需要忽略的签名字段
	Timeout       time.Duration // 签名过期时间
	Crypto        string        // 加密算法
}

//NewSigner 新签名
func NewSigner(opt Option) *Signer {
	// 解析参数
	opt = ParseCfg(opt)
	// 根据配置 赋值
	defaultKey := DefaultKeyName{
		keyNameTimestamp: opt.KeyNameTimestamp,
		keyNameNonceStr:  opt.KeyNameNonceStr,
		keyNameAppID:     opt.KeyNameAppID,
		keyNameSign:      opt.KeyNameSign,
		keyNameAppSecret: opt.KeyNameAppSecret,
	}
	signer := &Signer{
		DefaultKeyName: defaultKey,
		body:           make(url.Values),
		otherFields:    opt.OtherFields,
		skippedFields:  opt.SkippedFields,
		timeout:        opt.Timeout,
		secretVal:      opt.SecretVal,
		// bodyPrefix:     opt.bodyPrefix,
		// bodySuffix:     opt.bodySuffix,
	}

	return signer
}

// Sign 加签函数
func Sign(values url.Values) string {
	return SignWithOption(values, Option{})
}

// SignWithOption 加签函数
func SignWithOption(values url.Values, opt Option) string {
	// 创建签名对象
	v := NewSigner(opt)
	v.Init(values)

	str := v.GetSignature()

	return str
}

// Verify 验签函数
func Verify(values url.Values) error {
	return VerifyWithOption(values, Option{})
}

//VerifyWithOption 校验签名
func VerifyWithOption(values url.Values, opt Option) error {
	// 创建签名对象
	v := NewSigner(opt)
	v.Init(values)

	// 校验
	var err error
	// 验证字符串是否齐全
	// 校验时间
	err = v.CheckTimeStamp()
	if err != nil {
		return err
	}

	// 校验随机字符串

	// 校验签名
	err = v.CheckSignature()
	if err != nil {
		return err
	}
	// 校验通过
	return err
}

//Debug 检验
func Debug(values url.Values) map[string]string {
	return DebugWithOption(values, Option{})
}

//DebugWithOption 校验配置
func DebugWithOption(values url.Values, opt Option) map[string]string {
	// 创建签名对象
	v := NewSigner(opt)
	v.Init(values)
	return v.Debug()
}

// ParseCfg 解析配置
func ParseCfg(opt Option) Option {
	// 如果不存在，则赋值默认的
	if opt.KeyNameTimestamp == "" {
		opt.KeyNameTimestamp = KeyNameTimeStamp
	}
	if opt.KeyNameNonceStr == "" {
		opt.KeyNameNonceStr = KeyNameNonceStr
	}
	if opt.KeyNameAppID == "" {
		opt.KeyNameAppID = KeyNameAppID
	}
	if opt.KeyNameSign == "" {
		opt.KeyNameSign = KeyNameSign
	}
	if opt.KeyNameAppSecret == "" {
		opt.KeyNameAppSecret = KeyNameAppSecret
	}
	// if opt.bodyPrefix == "" {
	// 	opt.bodyPrefix = ""
	// }
	// if opt.bodySuffix == "" {
	// 	opt.bodySuffix = ""
	// }
	if opt.Timeout == 0 {
		opt.Timeout = 5 * time.Minute
	}
	return opt
}

//ParseJSON 解析json
func ParseJSON(buf []byte) (url.Values, error) {
	m := make(url.Values)
	// 解析为json，并且转换为字符串
	var f map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(buf))
	d.UseNumber()
	_ = d.Decode(&f)

	// 循环处理为values格式
	for key, val := range f {
		temp := String(val)
		m[key] = []string{temp}
	}

	return m, nil
}
