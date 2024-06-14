package ssh_term

import (
	"bufio"
	"errors"
	"go-protector/server/biz/model/dto"
	"go-protector/server/internal/ssh/ssh_con"
	"golang.org/x/crypto/ssh"
	"io"
)

type Terminal struct {
	ssoSessionId uint64
	sshCli       *ssh_con.Client
	sshSession   *ssh.Session
	in           io.WriteCloser
	out          *bufio.Reader
}

func NewTerminal(req *dto.ConnectBySessionReq, param *ssh_con.ConnectParam) (term *Terminal, err error) {
	var sshCli *ssh_con.Client
	var sshSession *ssh.Session
	defer func() {
		if err == nil {
			return
		}
		if term != nil {
			_ = term.Close()
		}

	}()
	sshCli, err = ssh_con.Connect(param)
	if err != nil {
		return
	}
	sshSession, err = sshCli.SSHClient.NewSession()
	if err != nil {
		return
	}

	// 输出
	var out *bufio.Reader
	var outPipe io.Reader
	if outPipe, err = sshSession.StdoutPipe(); err != nil {
		return
	}
	out = bufio.NewReader(outPipe)

	// 输入
	var inPipe io.WriteCloser
	if inPipe, err = sshSession.StdinPipe(); err != nil {
		return
	}
	term = &Terminal{
		ssoSessionId: req.Id,
		sshCli:       sshCli,
		sshSession:   sshSession,
		in:           inPipe,
		out:          out,
	}

	err = term.sshSession.RequestPty("xterm-256color", req.H, req.W, ssh.TerminalModes{
		ssh.ECHO:          1,     // 是否需要回显输入
		ssh.TTY_OP_ISPEED: 14400, // 速率
		ssh.TTY_OP_OSPEED: 14400, // 速率
	})
	if err = term.sshSession.Shell(); err != nil {
		return
	}

	return
}

// Write 写入
func (_self *Terminal) Write(in []byte) (int, error) {
	if _self.in == nil {
		return 0, errors.New("write pipe not available")
	}
	return _self.in.Write(in)
}

// WindowChange informs the remote host about a terminal window dimension change to h rows and w columns.
func (_self *Terminal) WindowChange(h, w int) error {
	return _self.sshSession.WindowChange(h, w)
}

// ReadRune 读取数据
func (_self *Terminal) ReadRune() (rune, int, error) {
	if _self.out == nil {
		return ' ', 0, errors.New("read pipe not available ")
	}
	return _self.out.ReadRune()
}

func (_self *Terminal) Close() error {

	if _self.in != nil {
		_ = _self.in.Close()
	}
	if _self.sshSession != nil {
		_ = _self.sshSession.Close()
	}
	if _self.sshCli != nil {
		_ = _self.sshCli.Close()
	}

	return nil
}
