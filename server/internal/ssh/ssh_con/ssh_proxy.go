package ssh_con

import (
	"context"
	"errors"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_logger"
	"go-protector/server/internal/database"
	"go-protector/server/internal/utils/async"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"strconv"
	"sync"
)

var sshProxyCache sync.Map

// 6-20	【实战】通过SSH网关连接资产，创建SSH网关连接-掌握sync.Map与并发控制；掌握DCL单例模式）
func GetSSHProxyById(id uint64) (sshProxy *SSHProxy, err error) {
	if id <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	defer func() {
		if err != nil {
			sshProxyCache.Delete(id)
		}
	}()

	value, _ := sshProxyCache.LoadOrStore(id, &SSHProxy{
		ID:         id,
		channelMap: make(map[string]*SSHChannel),
	})
	sshProxy = value.(*SSHProxy)
	// DCL
	if sshProxy.SSHClient != nil {
		return
	}

	sshProxy.proxyMutex.Lock()
	defer sshProxy.proxyMutex.Unlock()
	if sshProxy.SSHClient != nil {
		return
	}
	var connectDTO *ConnectParam
	if connectDTO, err = GenerateGatewayConnect(id); err != nil {
		return
	}

	var cli *Client
	cli, err = Connect(connectDTO)
	c_logger.Debug("Connect ip: %s, port: %d", connectDTO.IP, connectDTO.Port)
	if err != nil {
		c_logger.Error("Connect ip: %s, port: %d, err: %v", connectDTO.IP, connectDTO.Port, err)

		return
	}

	sshProxy.SSHClient = cli.SSHClient
	// 本机随机启动一个端口监听

	sshProxy.Listener, err = net.Listen("tcp", ":0")
	c_logger.Debug("启动监听: %s", sshProxy.Listener.Addr().String())
	if err != nil {
		c_logger.Debug("启动监听: %s, 失败: %v", sshProxy.Listener.Addr().String(), err)

		return
	}
	// 获取监听的端口
	sshProxy.LocalAddr = sshProxy.Listener.Addr().(*net.TCPAddr)

	// 主死从随
	sshProxy.Listen()
	return
}

type SSHProxy struct {
	// ID 网关ID
	ID uint64
	// proxyMutex 互斥锁
	proxyMutex sync.Mutex
	// LocalAddr 本地监听地址
	LocalAddr net.Addr
	// Listener 本地监听
	Listener net.Listener
	// SSHClient Client 代理
	SSHClient *ssh.Client
	once      sync.Once
	// channelMutex 读写锁
	channelMutex sync.RWMutex
	// channelMap 通道映射
	channelMap map[string]*SSHChannel
}

// Close 关闭连接,如指定了 chanKey 则关闭指定的channel,如未指定,关闭所有
func (_self *SSHProxy) Close(channelKey string) (err error) {
	var errs []error
	defer func() {
		if len(errs) > 0 {
			if err = errors.Join(errs...); err != nil {

				c_logger.Error("sshProxy err: %v", err)
			}

		}
	}()
	c_logger.Debug("sshProxy close channelKey: %s", channelKey)
	_self.proxyMutex.Lock()
	defer _self.proxyMutex.Unlock()
	_self.channelMutex.Lock()
	defer _self.channelMutex.Unlock()
	// 6-20	【实战】通过SSH网关连接资产，创建SSH网关连接-掌握sync.Map与并发控制；掌握DCL单例模式）
	if len(channelKey) > 0 {
		channel, ok := _self.channelMap[channelKey]
		if !ok {
			return
		}
		if channel.remoteConn != nil {

			errs = append(errs, channel.remoteConn.Close())
		}
		//errs = append(errs, channel.localConn.Close())
		delete(_self.channelMap, channelKey)
	} else {
		for k, channel := range _self.channelMap {
			errs = append(errs, channel.Close())
			delete(_self.channelMap, k)
		}
	}
	if len(_self.channelMap) > 0 {
		return
	}
	c_logger.Debug("没有在用 channel 开始关闭连接")
	errs = append(errs, _self.SSHClient.Close())
	if _self.Listener != nil {
		c_logger.Debug("关闭监听: %s", _self.Listener.Addr().String())
		errs = append(errs, _self.Listener.Close())
	}
	sshProxyCache.Delete(_self.ID)

	return
}

// Listen 启动本地监听,只启动一次
func (_self *SSHProxy) Listen() {
	_self.once.Do(func() {
		async.MainWork.AsyncRun(func() {

			for {
				conn, err := _self.Listener.Accept()
				if err != nil {
					return
				}
				go _self.handle(conn)

			}

		})
	})

}

// handle 处理新请求
func (_self *SSHProxy) handle(conn net.Conn) {
	var channel *SSHChannel
	defer func() {
		_ = conn.Close()
		if channel != nil {
			_ = channel.Close()
		}
	}()
	var err error
	var ok bool
	remoteAddr := conn.RemoteAddr()
	channelKey := remoteAddr.String()

	if channel, ok = _self.getSSHChannel(channelKey); !ok {
		return
	}

	//map 并发读写的问题
	//if channel, ok = _self.channelMap[channelKey]; !ok {
	//	return
	//}
	c_logger.Debug("sshLocal: %s, sshRemote: %s, channel: %s, targetAddr: %s",
		_self.SSHClient.LocalAddr().String(), _self.SSHClient.RemoteAddr().String(), channelKey, channel.targetAddr)
	// vi /etc/ssh/sshd_config  AllowTcpForwarding yes
	if channel.remoteConn, err = _self.SSHClient.Dial("tcp", channel.targetAddr); err != nil {
		c_logger.Error("dial: %s, err: %v", channel.targetAddr, err)
		//_, _ = conn.Write([]byte(err.Error()))
		return
	}
	channel.changeData(conn)

}

// createChannel 创建通道,保存本地连接信息,以及目标地址
func (_self *SSHProxy) createChannel(conn net.Conn, addr string) (channel *SSHChannel, err error) {

	var host, port string
	if host, port, err = net.SplitHostPort(addr); err != nil {
		return
	}
	tcpAddr := new(net.TCPAddr)
	tcpAddr.IP = net.ParseIP(host)

	if tcpAddr.Port, err = strconv.Atoi(port); err != nil {
		return
	}
	_self.channelMutex.Lock()
	defer _self.channelMutex.Unlock()
	channel = &SSHChannel{
		key:        conn.LocalAddr().String(),
		targetAddr: addr,
		remoteConn: nil,  // 等待连接成功后,这个是ssh提供的通道连接
		localConn:  conn, // 本地连接与本地监听代理的连接
		sshProxy:   _self,
	}
	c_logger.Debug("createChannel: %s", channel.key)
	_self.channelMap[channel.key] = channel
	return
}

// getSSHChannel 获取channel 读锁
// 6-21	【实战】创建本地监听与SSH网关通信-掌握GO语言读写锁
func (_self *SSHProxy) getSSHChannel(key string) (channel *SSHChannel, ok bool) {
	if len(key) <= 0 {
		return
	}
	c_logger.Debug("getSSHChannel: %s", key)
	_self.channelMutex.RLock()
	defer _self.channelMutex.RUnlock()
	channel, ok = _self.channelMap[key]
	return
}

type SSHChannel struct {
	key        string    // key
	targetAddr string    // 目标访问地址
	remoteConn net.Conn  // 远端连接
	localConn  net.Conn  // 本地连接
	sshProxy   *SSHProxy // ssh连接
}

// Close 关闭
func (_self *SSHChannel) Close() (err error) {
	return _self.sshProxy.Close(_self.key)
}

// changeData 交换数据
func (_self *SSHChannel) changeData(conn net.Conn) {

	go func() {
		_, err := io.Copy(_self.remoteConn, conn)
		if err != nil && !errors.Is(err, io.EOF) {
			c_logger.Error("Copy to remote error: %v", err)
		}
	}()
	_, _ = io.Copy(conn, _self.remoteConn)

}

var GenerateGatewayConnect = func(id uint64) (connectParam *ConnectParam, err error) {
	var model entity.AssetGateway
	db := database.GetDB(context.Background())
	if err = db.First(&model, id).Error; err != nil {
		return
	}
	connectParam = &ConnectParam{
		ID:       id,
		IP:       model.AgIp,
		Port:     model.AgPort,
		Username: model.AgAccount,
		Password: model.AgPassword,
	}
	return
}
