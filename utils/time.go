package utils

import "time"

// 时间操作

const (
	//SHORTTIMESTRING 短时间
	SHORTTIMESTRING = "20060102"
	//TIMEFORMAT 时间格式化
	TIMEFORMAT = "20060102150405"
	//NORMALTIMEFORMAT 普通格式化
	NORMALTIMEFORMAT = "2006-01-02 15:04:05"
)

//GetTime 当前时间
func GetTime() time.Time {
	return time.Now()
}

//GetTimeString 格式化为： 20060102150405
func GetTimeString(t time.Time) string {
	return t.Format(TIMEFORMAT)
}

//GetNormalTimeString 格式化为：2006-01-02 15:04:05
func GetNormalTimeString(t time.Time) string {
	return t.Format(NORMALTIMEFORMAT)
}

//GetTimeUnix 转为时间戳：秒数
func GetTimeUnix(t time.Time) int64 {
	return t.Unix()
}

//GetTimeMills 转为时间戳: 毫秒数
func GetTimeMills(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

//GetTimeByInt 时间戳转时间
func GetTimeByInt(t1 int64) time.Time {
	return time.Unix(t1, 0)
}

//GetTimeByString 字符串转时间
func GetTimeByString(timestring string) (time.Time, error) {
	if timestring == "" {
		return time.Time{}, nil
	}
	return time.ParseInLocation(TIMEFORMAT, timestring, time.Local)
}

//GetTimeByNormalString 标准字符串 转 时间
func GetTimeByNormalString(timestring string) (time.Time, error) {
	if timestring == "" {
		return time.Time{}, nil
	}
	return time.ParseInLocation(NORMALTIMEFORMAT, timestring, time.Local)
}

//CompareTime 比较两个时间大小
func CompareTime(t1, t2 time.Time) bool {
	return t1.Before(t2)
}

//GetNextHourTime n小时候后的时间字符串
func GetNextHourTime(s string, n int64) string {
	t2, _ := time.ParseInLocation(TIMEFORMAT, s, time.Local)
	t1 := t2.Add(time.Hour * time.Duration(n))
	return GetTimeString(t1)
}

//GetHourDiffer 计算两个时间差多少小时
func GetHourDiffer(startTime, endTime string) float32 {
	var hour float32
	t1, err := time.ParseInLocation(TIMEFORMAT, startTime, time.Local)
	t2, err := time.ParseInLocation(TIMEFORMAT, endTime, time.Local)
	if err == nil && CompareTime(t1, t2) {
		diff := GetTimeUnix(t2) - GetTimeUnix(t1)
		hour = float32(diff) / 3600
		return hour
	}
	return hour
}

//Checkhours 判断当前时间是否整点
func Checkhours() bool {
	_, m, s := GetTime().Clock()
	if m == s && m == 0 && s == 0 {
		return true
	}
	return false
}

//StringToNormalString 时间字符串转为标准字符串
func StringToNormalString(t string) string {
	if !(len(TIMEFORMAT) == len(t) || len(SHORTTIMESTRING) == len(t)) {
		return t
	}
	if len(SHORTTIMESTRING) == len(t) {
		t += "000000"
	}
	if len(TIMEFORMAT) == len(t) {
		t1, err := GetTimeByString(t)
		if err != nil {
			return t
		}
		t = GetTimeString(t1)
	}
	return t
}
