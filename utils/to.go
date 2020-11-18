package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

//JSONToMAP Convert json string to map
func JSONToMAP(jsonStr string) (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil, err
	}

	for k, v := range m {
		fmt.Printf("%v: %v\n", k, v)
	}

	return m, nil
}

//MAPToJSON Convert map json string
func MAPToJSON(m map[string]string) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return "", nil
	}

	return string(jsonByte), nil
}

/*
#string到int
int,err:=strconv.Atoi(string)
#string到int64
int64, err := strconv.ParseInt(string, 10, 64)
#int到string
string:=strconv.Itoa(int)
#int64到string
string:=strconv.FormatInt(int64,10)
*/

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
		return fmt.Sprintf("%v", val)
	}
}

//Int64 change val type to int64
func Int64(val interface{}) int64 {
	if val == nil {
		return 0
	}

	switch t := val.(type) {
	case bool:
		if t {
			return int64(1)
		}
		return int64(0)
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return int64(t)
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case float32:
		return int64(t)
	case float64:
		return int64(t)
	case []byte:
		i, _ := strconv.Atoi(string(t))
		return int64(i)
	case string:
		b, err := strconv.ParseBool(t)
		if err == nil {
			if b {
				return int64(1)
			}

			return int64(0)
		}
		i, _ := strconv.ParseFloat(t, 64)
		return int64(i)
	default:
		i, _ := strconv.ParseFloat((fmt.Sprintf("%v", t)), 64)
		return int64(i)
	}
}

//Float64 change val type to float64
func Float64(val interface{}) float64 {
	if val == nil {
		return float64(0)
	}

	switch t := val.(type) {
	case bool:
		if t {
			return float64(1)
		}

		return float64(0)
	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case uint:
		return float64(t)
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case float32:
		return float64(t)
	case float64:
		return t
	case []byte:
		i, _ := strconv.ParseFloat(string(t), 64)
		return i
	case string:
		i, _ := strconv.ParseFloat(t, 64)
		return i
	default:
		return float64(0)
	}
}

//Int change val type to int
func Int(val interface{}) int {
	return int(Int64(val))
}

//Int32 ToInt32 change val type to int32
func Int32(val interface{}) int32 {
	return int32(Int64(val))
}

//Float32 ToFloat32 change type to float32
func Float32(val interface{}) float32 {
	return float32(Float64(val))
}

//Decimal d
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
