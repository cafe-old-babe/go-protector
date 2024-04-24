package sshCli

import (
	"fmt"
	"go-protector/server/internal/custom/c_error"
	"golang.org/x/crypto/ssh"
	"reflect"
	"time"
)

// Connect 连接ssh
func Connect(dto *ConnectDTO) (cli *ssh.Client, err error) {
	if dto == nil || reflect.ValueOf(dto).IsZero() {
		err = c_error.ErrParamInvalid
		return
	}
	// 配置SSH连接
	timeout := dto.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	config := &ssh.ClientConfig{
		Timeout: timeout,
		User:    dto.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(dto.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到SSH主机
	addr := fmt.Sprintf("%s:%d", dto.IP, dto.Port)

	cli, err = ssh.Dial("tcp", addr, config)

	return
}
