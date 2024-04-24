package sshCli

import "time"

// ConnectDTO 连接对象
type ConnectDTO struct {
	ID       uint64
	IP       string
	Port     uint
	Username string
	Password string
	Timeout  time.Duration
}
