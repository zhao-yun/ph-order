package helper

import (
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return *(*string)(unsafe.Pointer(&b))
}

func S2I64(s string) int64 {
	i64, _ := strconv.ParseInt(s, 10, 64)
	return i64
}

func S2UI32(s string) uint32 {
	ui32, _ := strconv.ParseUint(s, 10, 32)
	return uint32(ui32)
}

func S2I32(s string) int32 {
	i32, _ := strconv.ParseInt(s, 10, 32)
	return int32(i32)
}

func I642S(i64 int64) string {
	s := strconv.FormatInt(i64, 10)
	return s
}

func I322S(i32 int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(i32)
	signed := i < 0
	if signed {
		i = -i
	}

	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func SSlice2I64Slice(ss []string) []int64 {
	if ss == nil {
		return nil
	}
	ii := make([]int64, 0, len(ss))
	for _, s := range ss {
		ii = append(ii, S2I64(s))
	}
	return ii
}

func I64Slice2SSlice(ii []int64) []string {
	if ii == nil {
		return nil
	}
	ss := make([]string, 0, len(ii))
	for _, i := range ii {
		ss = append(ss, I642S(i))
	}
	return ss
}

func I642ISlice(n int64) []int {
	result := make([]int, 0, 16)

	for n != 0 {
		result = append(result, int(n%10))
		n /= 10
	}

	length := len(result)
	for i := 0; i < length/2; i++ {
		result[i], result[length-1-i] = result[length-1-i], result[i]
	}
	return result
}

func S2I64Ptr(s string) *int64 {
	ret := S2I64(s)
	return &ret
}

func SPtr2I64Ptr(s *string) *int64 {
	if s == nil {
		return nil
	}
	ret := S2I64(*s)
	return &ret
}

func I642SPtr(i64 int64) *string {
	ret := I642S(i64)
	return &ret
}

func I64Ptr2SPtr(i64 *int64) *string {
	if i64 == nil {
		return nil
	}
	ret := I642S(*i64)
	return &ret
}

func Bytes2Str(b []byte) string {
	// convert SliceHeader to StringHeader
	return *(*string)(unsafe.Pointer(&b))
}

func Str2Bytes(s string) []byte {
	// convert StringHeader to SliceHeader
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	}
	return ""
}

func Bool2Str(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
