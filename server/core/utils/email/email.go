package email

import (
	"fmt"
	"go-protector/server/core/config"
	"go-protector/server/core/custom/c_error"
	"net/smtp"
)

const htmlFormat = `
<html>
	<body>
		<p>%s</p><br>
		<img src="%s" alt="图片">		
	</body>
</html>
`
const messageFormatText = "From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s"
const messageFormatHtml = "From: %s\r\nTo: %s\r\nSubject: %s\r\n%s"

type SendDTO struct {
	Email         *config.Email
	To            string
	Subject       string
	Body          string
	messageFormat string
}

// Send 发送邮件
func Send(dto *SendDTO) (err error) {
	if dto == nil || len(dto.To) <= 0 || len(dto.Subject) <= 0 || len(dto.Body) <= 0 {
		return c_error.ErrParamInvalid
	}
	var message string
	var host, port, from, password = dto.Email.Host, dto.Email.Port, dto.Email.Username, dto.Email.Password

	if len(host) <= 0 {
		host = config.GetConfig().Email.Host
	}
	if port <= 0 {
		port = config.GetConfig().Email.Port
	}
	if len(from) <= 0 {
		from = config.GetConfig().Email.Username
	}
	if len(password) <= 0 {
		password = config.GetConfig().Email.Password
	}
	if len(password) <= 0 || len(from) <= 0 || len(host) <= 0 || port <= 0 {
		return c_error.ErrParamInvalid
	}
	if len(dto.messageFormat) <= 0 {
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
func SendImage(dto SendDTO, imageBase64 string) (err error) {
	if len(dto.To) <= 0 || len(dto.Subject) <= 0 || len(dto.Body) <= 0 || len(imageBase64) <= 0 {
		return c_error.ErrParamInvalid
	}
	dto.Body = fmt.Sprintf("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		fmt.Sprintf(htmlFormat, dto.Body, imageBase64))
	dto.messageFormat = messageFormatHtml

	return Send(&dto)
}
