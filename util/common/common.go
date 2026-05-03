package common

import (
	"crypto/rand"
	"database/sql/driver"
	"fmt"
	"math/big"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`null`), nil
	}
	// 格式化为 "2006-01-02 15:04:05"
	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}
	str := string(data[1 : len(data)-1])
	parsedTime, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return fmt.Errorf("解析时间失败: %v", err)
	}
	t.Time = parsedTime
	return nil
}

func (t *Timestamp) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case []byte:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		t.Time = parsedTime
		return nil
	default:
		return fmt.Errorf("无法扫描类型 %T 到 Timestamp", value)
	}
}

func (t Timestamp) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Format("2006-01-02"))), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Time = time.Time{} // 处理 null
		return nil
	}
	str := string(data[1 : len(data)-1])
	parsedTime, err := time.Parse("2006-01-02", str)
	if err != nil {
		return fmt.Errorf("解析日期失败: %v (输入格式应为 YYYY-MM-DD)", err)
	}
	d.Time = parsedTime
	return nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		parsedTime, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = parsedTime
		return nil
	default:
		return fmt.Errorf("无法扫描类型 %T 到 Date", value)
	}
}

func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

func GenerateRandomCode(length int) string {
	code := ""
	max := big.NewInt(10) // 0-9的数字范围

	for i := 0; i < length; i++ {
		// 生成0-9的随机数
		randomNum, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "1547"
		}
		code += randomNum.String()
	}

	return code
}

func ParseTimestamp(timestamp string) (*Timestamp, error) {
	tmp, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		return nil, err
	}
	return &Timestamp{tmp}, nil
}
