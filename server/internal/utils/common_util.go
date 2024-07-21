package utils

// 生成5位随机数
import (
	"context"
	"fmt"
	"go-protector/server/internal/cache"
	"go-protector/server/internal/custom/c_type"
	"math"
	"math/rand"
	"strconv"
	"sync"
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

var GenerateDateNextSeq func(string) (string, error)

func init() {
	GenerateDateNextSeq = generateDateNextSeqFunc()

}

// generateDateNextSeqFunc 2023060800001 闭包+DCL
func generateDateNextSeqFunc() func(string) (string, error) {
	lastDateStr := time.Now().Add(-(24 * time.Hour)).Format("20060102")
	ctx := context.Background()
	fmt.Printf("init lastDateStr %s\n", lastDateStr)
	var lock sync.Mutex
	return func(pre string) (string, error) {

		newDateStr := time.Now().Format("20060102")

		redisClient := cache.GetRedisClient()
		var key = fmt.Sprintf("sequence:%s", lastDateStr)
		if newDateStr != lastDateStr {
			lock.Lock()
			defer lock.Unlock()
			if newDateStr != lastDateStr {
				lastDateStr = newDateStr
				key = fmt.Sprintf("sequence:%s", lastDateStr)
				nx := redisClient.SetNX(ctx, key, 0, 24*time.Hour+5*time.Minute)
				if nx.Err() != nil {
					_ = fmt.Errorf("connect failure: %+v", nx.Err())
					return "", nx.Err()
				}
			}

		}

		incr := redisClient.Incr(ctx, key)
		formattedNumber := fmt.Sprintf("%06d", incr.Val())

		return pre + lastDateStr + formattedNumber, nil
	}

}
