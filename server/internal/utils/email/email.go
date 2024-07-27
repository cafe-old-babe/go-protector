package email

import (
	"context"
	"errors"
	"fmt"
	"go-protector/server/internal/cache"
	"go-protector/server/internal/config"
	"go-protector/server/internal/custom/c_error"
	"net/smtp"
	"reflect"
	"strings"
	"time"
)

const htmlFormat = `
<html>
	<body>
		<p>%s</p><br>
		<img src="%s" alt="图片"/>		
	</body>
</html>
`
const messageFormatText = "From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s"
const messageFormatHtml = "From: %s\r\nTo: %s\r\nSubject: %s\r\n%s"

type SendDTO struct {
	Email         config.Email
	To            string
	Subject       string
	Body          string
	messageFormat string
}

// Send 发送邮件
// 5-8	【实战】邮箱验证码认证-掌握QQ邮箱发送认证码，防止恶意频繁发送验证码
func Send(dto *SendDTO) (err error) {
	if dto == nil || len(dto.To) <= 0 || len(dto.Subject) <= 0 || len(dto.Body) <= 0 {
		return c_error.ErrParamInvalid
	}
	var message string
	if reflect.ValueOf(dto.Email).IsZero() {
		dto.Email = config.GetConfig().Email
	}
	var host, port, from, password = dto.Email.Host, dto.Email.Port, dto.Email.Username, dto.Email.Password

	if len(password) <= 0 || len(from) <= 0 || len(host) <= 0 || port <= 0 {
		return c_error.ErrParamInvalid
	}
	if len(dto.messageFormat) <= 0 {
		// 5-9	【实战】OTP认证一-掌握使用QQ邮箱发送HTML格式的图片
		dto.messageFormat = messageFormatText
	}
	// 将邮件内容组织成RFC 822格式
	message = fmt.Sprintf(dto.messageFormat, from, dto.To, dto.Subject, dto.Body)
	// 连接SMTP服务器
	err = smtp.SendMail(fmt.Sprintf("%s:%d", host, port),
		smtp.PlainAuth("", from, password, host),
		from, []string{dto.To}, []byte(message))

	return
}

// SendImage 发送带图片
// 5-9	【实战】OTP认证一-掌握使用QQ邮箱发送HTML格式的图片
func SendImage(dto SendDTO, imageBase64 string) (err error) {
	if len(dto.To) <= 0 || len(dto.Subject) <= 0 || len(dto.Body) <= 0 || len(imageBase64) <= 0 {
		return c_error.ErrParamInvalid
	}
	if !strings.HasPrefix(imageBase64, "data:image/png;base64,") {
		imageBase64 = "data:image/png;base64," + imageBase64
	}

	dto.Body = fmt.Sprintf("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		fmt.Sprintf(htmlFormat, dto.Body, imageBase64))
	dto.messageFormat = messageFormatHtml

	return Send(&dto)
}

// VerifySendInterval 校验发送间隔时间
func VerifySendInterval(context context.Context, key string,
	expireTime, interval time.Duration) (err error) {
	if interval.Seconds() >= expireTime.Seconds() {
		err = c_error.ErrParamInvalid
		return err
	}

	redisClient := cache.GetRedisClient()
	var ttl time.Duration
	if ttl, err = redisClient.TTL(context, key).Result(); err != nil {
		return err
	}
	ttlSeconds := ttl.Seconds()
	if ttlSeconds <= 0 ||
		ttlSeconds < expireTime.Seconds()-interval.Seconds() {
		return nil
	}
	return errors.New("发送过于频繁,请稍后再试")

}
