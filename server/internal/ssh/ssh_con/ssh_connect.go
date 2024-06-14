package ssh_con

import (
	"errors"
	"fmt"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_logger"
	"golang.org/x/crypto/ssh"
	"net"
	"reflect"
	"time"
)

type Client struct {
	SSHClient  *ssh.Client
	SSHChannel *SSHChannel
}

// ConnectParam 连接对象
type ConnectParam struct {
	ID        uint64
	IP        string
	Port      int
	Username  string
	Password  string
	Timeout   time.Duration
	GatewayId uint64
}

// Connect 连接ssh
func Connect(param *ConnectParam) (cli *Client, err error) {
	if param == nil || reflect.ValueOf(param).IsZero() {
		err = c_error.ErrParamInvalid
		return
	}
	// 配置SSH连接
	timeout := param.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	config := &ssh.ClientConfig{
		Timeout: timeout,
		User:    param.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(param.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var sshProxy *SSHProxy
	if param.GatewayId > 0 {
		if sshProxy, err = GetSSHProxyById(param.GatewayId); err != nil {
			return
		}
	}

	// 连接到SSH主机
	addr := fmt.Sprintf("%s:%d", param.IP, param.Port)
	cli = new(Client)
	if sshProxy == nil {
		cli.SSHClient, err = ssh.Dial("tcp", addr, config)
		return
	}
	// 通过网关连接 必须加锁
	sshProxy.proxyMutex.Lock()
	defer sshProxy.proxyMutex.Unlock()

	var conn net.Conn
	// 连接本地监听端口
	localListenerAddr := sshProxy.Listener.Addr().String()
	c_logger.Debug("connect sshProxy listen: %s, target: %s", localListenerAddr, addr)
	if conn, err = net.DialTimeout("tcp",
		localListenerAddr, config.Timeout); err != nil {
		c_logger.Error("connect sshProxy listen: %s, err: %v ", localListenerAddr, err)
		return
	}

	if cli.SSHChannel, err = sshProxy.createChannel(conn, addr); err != nil {
		c_logger.Error("createChannel err: %v ", localListenerAddr, err)
		return

	}

	// 此处开始握手
	var c ssh.Conn
	var chans <-chan ssh.NewChannel
	var reqs <-chan *ssh.Request
	if c, chans, reqs, err = ssh.NewClientConn(conn, localListenerAddr, config); err != nil {
		c_logger.Error("NewClientConn err: %v ", err)
		return
	}
	cli.SSHClient = ssh.NewClient(c, chans, reqs)
	return

}

func (_self *Client) Close() (err error) {
	var errs []error
	if _self.SSHClient != nil {
		errs = append(errs, _self.SSHClient.Close())
	}
	if _self.SSHChannel != nil {
		errs = append(errs, _self.SSHChannel.Close())
	}

	return errors.Join(errs...)
}
