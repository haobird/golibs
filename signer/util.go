package signer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var randomSrc = rand.NewSource(time.Now().UnixNano())

// RandString 返回指定长度的随机字符串。
// 最小长度为4，最大长度为1024
func RandString(num int) string {
	if num < 4 {
		num = 4
	}
	if num > 1024 {
		num = 1024
	}
	bytes := make([]byte, num)
	for i, cache, remain := num-1, randomSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			bytes[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(bytes)
}

//md5V md5
func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func encrypt(str string) string {
	str = md5V(str)
	return strings.ToUpper(str)
}

// SortKVPairs 将Map的键值对，按字典顺序拼接成字符串
func SortKVPairs(m url.Values) string {
	size := len(m)
	if size == 0 {
		return ""
	}
	keys := make([]string, size)
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	pairs := make([]string, size)
	for i, key := range keys {
		pairs[i] = key + "=" + strings.Join(m[key], ",")
	}
	return strings.Join(pairs, "&")
}

//String change val type to string
func String(val interface{}) string {
	if val == nil {
		return ""
	}

	switch t := val.(type) {
	case bool:
		return strconv.FormatBool(t)
	case int:
		return strconv.FormatInt(int64(t), 10)
	case int8:
		return strconv.FormatInt(int64(t), 10)
	case int16:
		return strconv.FormatInt(int64(t), 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case []byte:
		return string(t)
	case string:
		return t
	default:
		data, _ := json.Marshal(val)
		return string(data)
	}
}
