package utils

// 生成5位随机数
import (
	"fmt"
	"go-protector/server/internal/custom/c_type"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func GetRandomCode() string {
	return GetRandomCodeI(0)
}

func GetRandomCodeI(digits int) string {
	if digits <= 0 {
		digits = 6
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	format := "%0" + strconv.Itoa(digits) + "d"
	pow := math.Pow(10, float64(digits))
	return fmt.Sprintf(format, r.Intn(int(pow)))
}

// ParseTime 转换
func ParseTime(str string) (time c_type.Time, err error) {
	err = time.UnmarshalJSON([]byte(str))
	return
}
