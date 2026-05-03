package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	Marshal    = json.Marshal
	Unmarshal  = json.Unmarshal
	NewEncoder = json.NewEncoder
	NewDecoder = json.NewDecoder
)

func ToJSON(i interface{}) string {
	data, err := json.Marshal(i)
	if err != nil {
		logrus.Warnf("json.Marshal failed, err: %v", err)
	}
	return string(data)
}

func ToJSONBytes(i interface{}) []byte {
	data, err := json.Marshal(i)
	if err != nil {
		logrus.Warnf("json.Marshal failed, err: %v", err)
	}
	return data
}

// ParseJSON 反序列.
func ParseJSON(str string, i interface{}) error {
	err := json.Unmarshal([]byte(str), i)
	if err != nil {
		logrus.Warnf("json.Unmarshal failed, str: %v, err: %v", str, err)
	}
	return err
}
