package utils

// 生成5位随机数
import (
	"fmt"
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
